package ihttp

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Rule struct of rule
type Rule struct {
	Key     string
	Type    string
	Rule    string
	Msg     string
	Default interface{}
	Enum    []interface{}
}

// struct of result for check
type res struct {
	Key   string
	Value interface{}
	Err   string
}

// ret map without sync
type ret map[string]interface{}

// Int get param of request for int
func (r ret) Int(k string) int {
	item, itemOk := r[k]
	if !itemOk {
		return 0
	}

	if v, ok := item.(int); ok {
		return v
	}

	return 0
}

// Int64 get param of request for int64
func (r ret) Int64(k string) int64 {
	item, itemOk := r[k]
	if !itemOk {
		return 0
	}

	if v, ok := item.(int64); ok {
		return v
	}

	return 0
}

// String get param of request for string
func (r ret) String(k string) string {
	item, itemOk := r[k]
	if !itemOk {
		return ``
	}

	if v, ok := item.(string); ok {
		return v
	}

	return ``
}

// Bool get param of request for bool
func (r ret) Bool(k string) bool {
	item, itemOk := r[k]
	if !itemOk {
		return false
	}

	if v, ok := item.(bool); ok {
		return v
	}

	return false
}

// Float64 get param of request for float64
func (r ret) Float64(k string) float64 {
	item, itemOk := r[k]
	if !itemOk {
		return 0
	}

	if v, ok := item.(float64); ok {
		return v
	}

	return 0
}

// Validate do validate
func Validate(r *http.Request, rules []Rule) (*ret, error) {
	ret, dataRes := make(ret), make([]res, len(rules))
	for i, rule := range rules {
		if strings.EqualFold(rule.Key, "") {
			dataRes[i] = res{}
			continue
		}

		key, value := getValue(rule.Key, r)
		if strings.EqualFold(value.(string), "") {
			if rule.Default != nil {
				value = rule.Default
			}
		} else {
			value = switchValue(value.(string), rule.Type)
		}

		if strings.EqualFold(rule.Rule, "") {
			dataRes[i] = res{key, value, ""}
			continue
		}
		if checkRule(value, rule) {
			dataRes[i] = res{key, value, ""}
			continue
		}

		if !strings.EqualFold(rule.Msg, "") {
			dataRes[i] = res{key, "", rule.Msg}
			continue
		}
		dataRes[i] = res{key, "", fmt.Sprintf("Parameter [%s] validate error", key)}
	}

	var errMsg []string
	for _, resItem := range dataRes {
		if !strings.EqualFold(resItem.Err, "") {
			errMsg = append(errMsg, resItem.Err)
			continue
		}
		if strings.EqualFold(resItem.Key, "") {
			continue
		}
		ret[resItem.Key] = resItem.Value
	}

	if errMsg != nil {
		return nil, errors.New(strings.Join(errMsg, "; "))
	}

	return &ret, nil
}

func getValue(key string, r *http.Request) (string, interface{}) {
	keyArr := strings.Split(key, "|")
	if len(keyArr) == 1 {
		return key, r.FormValue(key)
	}
	if strings.EqualFold(keyArr[1], "header") {
		return keyArr[0], r.Header.Get(keyArr[0])
	}
	return keyArr[0], ""
}

func switchValue(value string, typ string) interface{} {
	switch strings.ToLower(typ) {
	case "int":
		v, _ := strconv.Atoi(value)
		return v
	case `int64`:
		v, _ := strconv.ParseInt(value, 10, 64)
		return v
	case "bool":
		v, _ := strconv.ParseBool(value)
		return v
	case "float64":
		v, _ := strconv.ParseFloat(value, 64)
		return v
	case "string":
		return value
	default:
		return value
	}
}

func checkRule(value interface{}, rule Rule) bool {
	r := strings.Split(rule.Rule, "|")
	for _, i := range r {
		switch i {
		case "":
			continue
		case "required":
			if value == nil {
				return false
			}
			if v, ok := value.(string); ok && strings.EqualFold(v, "") {
				return false
			}
			continue
		case "include":
			flag := false
			for _, v := range rule.Enum {
				if value == v {
					flag = true
				}
			}
			if flag {
				continue
			}
			return false
		case "exclude":
			for _, v := range rule.Enum {
				if value == v {
					return false
				}
			}
			continue
		default:
			if strings.Contains(i, "between:") {
				return between(value, i)
			}
			return false
		}
	}
	return true
}

func between(v interface{}, r string) bool {
	if _, ok := v.(string); !ok {
		return false
	}
	var min, max int
	if _, err := fmt.Sscanf(r, "between:%d,%d", &min, &max); err != nil {
		return false
	}
	l := len(v.(string))
	return (min <= 0 || l >= min) && (max <= 0 || l <= max)
}
