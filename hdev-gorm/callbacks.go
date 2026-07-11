package hdev_gorm

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"net"
	"os"
	"strings"
	"syscall"
)

func registerRetryCallbacks(db *gorm.DB) {
	db.Callback().Query().Replace("gorm:query", retryQueryOnError)
	db.Callback().Row().Replace("gorm:row", retryRowOnError)
	db.Callback().Raw().Replace("gorm:raw", retryRawOnError)
}

func retryQueryOnError(db *gorm.DB) {
	callbacks.Query(db)
	if !needsRetry(db) {
		return
	}

	doRetrySecondary(db, callbacks.Query)
}

func retryRowOnError(db *gorm.DB) {
	callbacks.RowQuery(db)
	if !needsRetry(db) {
		return
	}

	doRetrySecondary(db, callbacks.RowQuery)
}

func retryRawOnError(db *gorm.DB) {
	callbacks.RawExec(db)
	if !needsRetry(db) {
		return
	}

	doRetryGuess(db, callbacks.RawExec)
}

func doRetrySecondary(db *gorm.DB, fn func(db *gorm.DB)) {
	if rawSQL := db.Statement.SQL.String(); len(rawSQL) > 0 {
		doRetryGuess(db, fn)
	} else {
		_, locking := db.Statement.Clauses["FOR"]
		if !locking {
			db.Statement.ConnPool = getSecondaryDb().ConnPool
			db.Error = nil
			fn(db)
		}
	}
}

func doRetryGuess(db *gorm.DB, fn func(db *gorm.DB)) {
	if rawSQL := strings.TrimSpace(db.Statement.SQL.String()); len(rawSQL) > 10 && strings.EqualFold(rawSQL[:6], "select") && !strings.EqualFold(rawSQL[len(rawSQL)-10:], "for update") {
		db.Statement.ConnPool = getSecondaryDb().ConnPool
		db.Error = nil
		fn(db)
	}
}

func needsRetry(db *gorm.DB) bool {
	if db.Error == nil {
		return false
	}

	//如果是带事务的，不要重试
	//因为事务可能包含多个执行，即使单个重试成功了，也不起作用
	if isTransaction(db.ConnPool) {
		return false
	}

	//没有备份库配置
	if getSecondaryDb() == nil {
		return false
	}

	return isNetworkError(db.Error) || isTimeoutError(db.Error)
}

func isNetworkError(err error) bool {
	if err == nil {
		return false
	}

	var opErr *net.OpError
	ok := errors.As(err, &opErr)
	if !ok {
		return false
	}

	var t *os.SyscallError
	switch {
	case errors.As(opErr.Err, &t):
		var errno syscall.Errno
		if errors.As(t.Err, &errno) {
			switch {
			case errors.Is(errno, syscall.ECONNREFUSED):
				//TCP 无法建立连接
				//connection refused
				return true
			case errors.Is(errno, syscall.ETIMEDOUT):
				//读写超时为例
				return true
			}
		}
	}

	return false
}

func isTimeoutError(err error) bool {
	if err == nil {
		return false
	}

	var netErr net.Error
	ok := errors.As(err, &netErr)

	if ok && netErr.Timeout() {
		return true
	}

	if strings.Contains(err.Error(), "timeout") {
		return true
	}
	return false
}

func isTransaction(connPool gorm.ConnPool) bool {
	_, ok := connPool.(gorm.TxCommitter)
	return ok
}
