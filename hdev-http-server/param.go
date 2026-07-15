package ihttp

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidKind  = fmt.Errorf("invalid kind")
	ErrNotCanSet    = fmt.Errorf("can not set")
	ErrInvalidParam = fmt.Errorf("invalid params")
)

// RequestParam 获取参数
func RequestParam(data interface{}, isReenter, isValidate bool, r *http.Request) error {
	body := r.Body
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}
	if isReenter {
		defer func() {
			r.Body = ioutil.NopCloser(bytes.NewBuffer(b))
		}()
	}
	switch r.Method {
	case "GET":
		return parseForm(data, r.URL.Query())
	case "POST", "PUT", "PATCH", "DELETE":
		switch r.Header.Get("Content-Type") {
		case "application/x-www-form-urlencoded":
			return parseForm(data, r.Form)
		case "application/json":
			return parseBody(data, isValidate, b)
		case "application/json; charset=utf-8":
			return parseBody(data, isValidate, b)
		case "text/plain":
			return parseBody(data, isValidate, b)
		default:
			return parseForm(data, r.PostForm)
		}
	}

	return nil
}

//// Validate 规则校验
//func Validate(obj interface{}) error {
//	return validate(obj)
//}

func parseBody(data interface{}, isValidate bool, b []byte) error {

	if err := json.Unmarshal(b, data); err != nil {
		return err
	}
	if isValidate {
		var valid *validator.Validate
		valid = validator.New()
		err := valid.Struct(data)
		return err
	}
	return nil
	//return validate(data)
}

func parseForm(data interface{}, value url.Values) error {

	refV := reflect.Indirect(reflect.ValueOf(data))
	refT := reflect.TypeOf(data)
	if refT.Kind() == reflect.Ptr {
		refT = refT.Elem()
	}
	if refV.Kind() != reflect.Struct {
		return ErrInvalidKind
	}

	if !refV.CanSet() {
		return ErrNotCanSet
	}
	for i := 0; i < refV.NumField(); i++ {
		t := refT.Field(i)
		v := refV.Field(i)
		tag := t.Tag.Get("json")
		if len(tag) == 0 {
			continue
		}

		form := value.Get(tag)
		if len(form) == 0 {
			continue
		}

		if err := setValue(v, form); err != nil {
			return err
		}
	}

	return validate(data)
}

func setValue(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.Int:
		return setIntField(value, 0, v)
	case reflect.Int8:
		return setIntField(value, 8, v)
	case reflect.Int16:
		return setIntField(value, 16, v)
	case reflect.Int32:
		return setIntField(value, 32, v)
	case reflect.Int64:
		return setIntField(value, 64, v)
	case reflect.Uint:
		return setUintField(value, 0, v)
	case reflect.Uint8:
		return setUintField(value, 8, v)
	case reflect.Uint16:
		return setUintField(value, 16, v)
	case reflect.Uint32:
		return setUintField(value, 32, v)
	case reflect.Uint64:
		return setUintField(value, 64, v)
	case reflect.String:
		v.SetString(value)
	case reflect.Bool:
		return setBoolField(value, v)
	case reflect.Float32:
		return setFloatField(value, 32, v)
	case reflect.Float64:
		return setFloatField(value, 64, v)
	case reflect.Struct:
		switch v.Interface().(type) {
		case time.Time:
			return nil
		}
		return json.Unmarshal([]byte(value), v.Addr().Interface())
	case reflect.Map:
		return json.Unmarshal([]byte(value), v.Addr().Interface())
	}

	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return err
}

func validate(obj interface{}) error {
	refV := reflect.Indirect(reflect.ValueOf(obj))
	if !refV.CanSet() {
		return ErrNotCanSet
	}

	return validateField(refV, "")
}

func validateField(val reflect.Value, rule string) error {
	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < val.Len(); i++ {
			if err := validateField(val.Index(i), ""); err != nil {
				return err
			}
		}
		return nil
	case reflect.Struct:
		for i := 0; i < val.NumField(); i++ {
			if err := validateField(val.Field(i), val.Type().Field(i).Tag.Get("rule")); err != nil {
				return err
			}
		}
	case reflect.Ptr:
		return validateField(val.Elem(), "")
	case reflect.Map, reflect.Chan, reflect.Func:
	default:
		if len(rule) == 0 {
			return nil
		}
		for _, v := range strings.Split(rule, "|") {
			switch v {
			case "required":
				if val.IsZero() {
					return ErrInvalidParam
				}
				// other rule
			}
		}
	}
	return nil
}
