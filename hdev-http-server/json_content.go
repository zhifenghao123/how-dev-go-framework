package ihttp

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

// JsonData get param of application/json
func JsonData(r *http.Request, i interface{}) error {
	dat, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.New(`illegal request: read body fail`)
	}

	if err := j.Unmarshal(dat, i); err != nil {
		return errors.New(`illegal request: not a json`)
	}

	return nil
}

// CheckJson check upsert params
func CheckJson(r *http.Request, data interface{}) error {
	jsonData, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return errors.New(`illegal request: read body fail`)
	}
	unmarshalErr := json.Unmarshal(jsonData, data)
	if unmarshalErr != nil {
		return errors.New(`illegal request: not a json`)
	}

	objType := reflect.TypeOf(data).Elem()
	names := make([]string, objType.NumField())
	for i := range names {
		names[i] = objType.Field(i).Name
	}

	var errMsg []string
	values := reflect.ValueOf(data).Elem()
	for i, name := range names {
		value := values.FieldByName(name)
		tag := objType.Field(i).Tag
		// 根据 tag 获取 json key 值
		jsonKey := tag.Get(`json`)
		valueStr := value.String()
		if jsonKey == `-` {
			continue
		}

		keyNote := tag.Get(`note`)
		if keyNote == `` {
			keyNote = jsonKey
		}

		// check empty
		if required := tag.Get(`required`); required == `1` && value.IsZero() {
			errMsg = append(errMsg, keyNote+` 字段不可为空`)
		}
		// check length
		if maxLen, lenRes := checkLen(tag, valueStr); !lenRes {
			errMsg = append(errMsg, keyNote+`字段长度超过`+maxLen)
		}

		// check type
		if typeRes, typeMsg := checkType(tag, valueStr); !typeRes {
			errMsg = append(errMsg, keyNote+typeMsg)
		}

		// check limit
		if limitRes, limitMsg := checkObjLimit(tag, valueStr); !limitRes {
			errMsg = append(errMsg, keyNote+limitMsg)
		}

		// check enum
		if enumRes, enumMsg := checkObjEnum(tag, valueStr); !enumRes {
			errMsg = append(errMsg, keyNote+enumMsg)
		}
	}

	if len(errMsg) != 0 {
		allErrMsg := strings.Join(errMsg, `;`)
		return errors.New(allErrMsg)
	}

	return nil
}

// checkLen check field length
func checkLen(tag reflect.StructTag, valueStr string) (maxLenStr string, checkRes bool) {
	maxLenStr = tag.Get(`maxLen`)
	if length, _ := strconv.Atoi(maxLenStr); length != 0 {
		paramLen := utf8.RuneCountInString(valueStr)
		if paramLen > length {
			return maxLenStr, false
		}
	}
	return maxLenStr, true
}

// checkType check field type
func checkType(tag reflect.StructTag, valueStr string) (bool, string) {
	fieldType := tag.Get(`type`)
	if valueStr == `` || fieldType == `` {
		return true, ``
	}
	switch fieldType {
	case `time`:
		if !checkTimeStr(valueStr) {
			return false, `不是时间戳`
		}
	case `image`:
		fileExt := strings.ToLower(filepath.Ext(valueStr))
		if fileExt == "" || !strings.Contains(`.png,.jpg,.jpeg,.gif`, fileExt) {
			return false, `不是图片`
		}
	default:
		return false, `错误的规则类型`

	}
	return true, ``
}

// checkTimeStr check time string
func checkTimeStr(timeStr string) bool {
	// 必须校验需用required另外限制
	if timeStr == `` {
		return true
	}

	_, err := time.Parse("2006-01-02 15:04:05", timeStr)
	return err == nil
}

// checkObjLimit check object limit
func checkObjLimit(tag reflect.StructTag, valueStr string) (bool, string) {
	limit := tag.Get(`include`)
	if limit == `` {
		return true, ``
	}
	limit = `,` + limit + `,`
	valueStr = `,` + valueStr + `,`
	if !strings.Contains(limit, valueStr) {
		return false, `不在限制内`
	}

	return true, ``
}

// checkObjEnum check object enum
func checkObjEnum(tag reflect.StructTag, valueStr string) (bool, string) {
	if valueStr == `` {
		return true, ``
	}
	enumLimit := tag.Get(`enum`)
	if enumLimit == `` {
		return true, ``
	}
	valueSet := strings.Split(valueStr, `,`)
	enumLimit = `,` + enumLimit + `,`
	for _, valueItem := range valueSet {
		valueItem = `,` + valueItem + `,`
		if !strings.Contains(enumLimit, valueItem) {
			return false, `不再枚举限制内`
		}
	}
	return true, ``
}
