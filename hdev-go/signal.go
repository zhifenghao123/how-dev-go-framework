package hdevgo

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

const pidFile = "pid"

// logPid 写入 pid 文件
func logPid() {
	pid := os.Getpid()
	log.Printf("[hdev-context] pid [%d]", pid)
	f, err := os.OpenFile(pidFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Printf("[hdev-context] pid file open error: %v", err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(strconv.Itoa(pid)); err != nil {
		log.Printf("[hdev-context] pid write error: %v", err)
	}
}

// removePid 删除 pid 文件
func removePid() error {
	return os.Remove(pidFile)
}

// readPid 读取 pid 文件
func readPid() (int, error) {
	r, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return -1, err
	}
	return strconv.Atoi(string(r))
}

// stopRunningProcess 给当前 pid 文件中的进程发送 SIGTERM
func stopRunningProcess() error {
	pid, err := readPid()
	if err != nil {
		return err
	}
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return process.Signal(syscall.SIGTERM)
}

// listenStopSignal 监听 SIGINT/SIGTERM，收到后向 stop 通道发送 true
func listenStopSignal(stop chan<- bool) {
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	for {
		msg := <-s
		switch msg {
		case syscall.SIGINT, syscall.SIGTERM:
			log.Printf("[hdev-context] received signal [%s]", msg)
			if stop != nil {
				stop <- true
			}
			return
		}
	}
}
