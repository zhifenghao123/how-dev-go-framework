package hdevconfig

import (
	"fmt"
	"github.com/goinggo/mapstructure"
	"log"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

type ValueExtractor struct {
	Key      string
	Type     string
	Value    interface{}
	Child    map[string]*ValueExtractor
	Callback []func(i interface{})
	sm       sync.RWMutex
}

func (v *ValueExtractor) Get(key string, b ...func(i interface{})) interface{} {
	d := v.value(strings.Split(key, "."))
	if d == nil {
		return nil
	}

	if b != nil {
		go callback(d, b...)
	}

	return d.Value
}

func (v *ValueExtractor) GetAsString(key string, b ...func(i interface{})) string {
	i := v.Get(key, b...)
	if i == nil {
		return ""
	}
	switch i.(type) {
	case string:
		return i.(string)
	case int:
		return strconv.Itoa(i.(int))
	case int64:
		return strconv.FormatInt(i.(int64), 10)
	default:
		return ""
	}
}

// Int get config by key, return value type of int
func (v *ValueExtractor) GetAsInt(key string, b ...func(i interface{})) int {
	i := v.Get(key, b...)
	if i == nil {
		return 0
	}
	switch i.(type) {
	case int64:
		return int(i.(int64))
	case int:
		return i.(int)
	case float64:
		return int(i.(float64))
	case bool:
		if i.(bool) {
			return 1
		}
		return 0
	case string:
		ii, _ := strconv.Atoi(i.(string))
		return ii
	default:
		return 0
	}
}

func (v *ValueExtractor) GetAsFloat(key string, b ...func(i interface{})) float64 {
	i := v.Get(key, b...)
	if i == nil {
		return 0
	}
	switch i.(type) {
	case float64:
		return i.(float64)
	case int:
		return float64(i.(int))
	case int64:
		return float64(i.(int64))
	case bool:
		if i.(bool) {
			return 1
		}
		return 0
	case string:
		ii, _ := strconv.ParseFloat(i.(string), 64)
		return ii
	default:
		return 0
	}
}

func (v *ValueExtractor) GetAsBool(key string, b ...func(i interface{})) bool {
	i := v.Get(key, b...)
	if i == nil {
		return false
	}
	switch i.(type) {
	case bool:
		return i.(bool)
	case int:
		return i.(int) != 0
	case int64:
		return i.(int64) != 0
	case float64:
		return i.(float64) != 0
	case string:
		return i.(string) != "0" && i.(string) != ""
	default:
		return false
	}
}

func (v *ValueExtractor) GetAsArray(key string, b ...func(i interface{})) []interface{} {
	i := v.Get(key, b...)
	if i == nil {
		return nil
	}
	if a, ok := i.([]interface{}); ok {
		return a
	}
	return nil
}

func (v *ValueExtractor) GetAsMap(key string, b ...func(i interface{})) map[string]interface{} {
	i := v.Get(key, b...)
	if i == nil {
		return nil
	}
	if mv, ok := i.(map[string]interface{}); ok {
		m := make(map[string]interface{})
		for k, item := range mv {
			m[k] = item
		}
		return m
	}
	return nil
}

func (v *ValueExtractor) GetAsStruct(key string, s interface{}, b ...func(i interface{})) interface{} {
	f := func(v interface{}) interface{} {
		if v == nil {
			return nil
		}
		ins, err := decode(v, reflect.TypeOf(s))
		if err != nil {
			log.Println("struct config error : ", err)
		}
		return ins
	}

	// format callback function
	for x, y := range b {
		t := y
		b[x] = func(v interface{}) {
			t(f(v))
		}
	}

	return f(v.Get(key, b...))
}

func (v *ValueExtractor) GetAsStructArray(key string, s interface{}, b ...func(i interface{})) []interface{} {
	f := func(v interface{}) []interface{} {
		if v == nil {
			return nil
		}

		var list []interface{}
		ty := reflect.TypeOf(s)
		if j, ok := v.([]interface{}); ok {
			for _, item := range j {
				ins, err := decode(item, ty)
				if err != nil {
					log.Println("struct array config error : ", err)
					continue
				}
				list = append(list, ins)
			}
			return list
		}

		if j, ok := v.([]map[string]interface{}); ok {
			for _, item := range j {
				ins, err := decode(item, ty)
				if err != nil {
					log.Println("struct array config error : ", err)
					continue
				}
				list = append(list, ins)
			}
			return list
		}

		return nil
	}

	// format callback function
	for x, y := range b {
		t := y
		b[x] = func(v interface{}) {
			t(f(v))
		}
	}

	return f(v.Get(key, b...))
}

func (v *ValueExtractor) Load(m interface{}, name string) {
	v.recursionValue(v, name, m)
}

func (v *ValueExtractor) value(key []string) *ValueExtractor {
	if c, ok := v.Child[key[0]]; ok {
		if len(key) == 1 {
			return c
		}
		return c.value(key[1:])
	}
	return nil
}

func (v *ValueExtractor) recursionValue(parent *ValueExtractor, k string, i interface{}) {
	var value *ValueExtractor
	key := k

	if parent.Key != "" {
		key = fmt.Sprintf("%s.%s", parent.Key, key)
	}

	var t string
	switch i.(type) {
	case int:
		t = _typeInt
		break
	case string:
		t = _typeString
		break
	case bool:
		t = _typeBoolean
		break
	case []interface{}:
		t = _typeArray
		break
	case map[string]interface{}:
		t = _typeMap
		break
	default:
		t = _typeUnknown
	}

	if d, ok := parent.Child[k]; ok {
		value = d
		if !reflect.DeepEqual(d.Value, i) {
			defer trigger(d)
		}
		value.Value = i
	} else {
		value = &ValueExtractor{
			Key:      key,
			Type:     t,
			Value:    i,
			Child:    make(map[string]*ValueExtractor),
			Callback: make([]func(i interface{}), 0),
		}
	}

	if t == _typeMap {
		for k1, v1 := range i.(map[string]interface{}) {
			v.recursionValue(value, k1, v1)
		}
	}

	parent.Child[k] = value
}

func trigger(value *ValueExtractor) {
	for _, c := range value.Callback {
		c(value.Value)
	}
}

// format callback function
func callback(value *ValueExtractor, b ...func(i interface{})) {
	value.sm.Lock()
	defer value.sm.Unlock()
	value.Callback = func(b []func(i interface{})) []func(i interface{}) {
		callback := value.Callback
		identify := true
		for _, c := range b {
			//while the first callback function is nil
			//make identity false and skip check whether the callback function is the same
			if reflect.ValueOf(c).IsNil() {
				identify = false
				continue
			}

			//check and skip append while the callback function is the same
			if identify && !func(b func(i interface{})) bool {
				for _, c := range callback {
					if reflect.ValueOf(b) == reflect.ValueOf(c) ||
						reflect.ValueOf(b).Pointer() == reflect.ValueOf(c).Pointer() {
						return false
					}
				}
				return true
			}(c) {
				continue
			}
			callback = append(callback, c)
		}
		return callback
	}(b)
}

func decode(v interface{}, t reflect.Type) (interface{}, error) {
	ins := reflect.New(t).Interface()
	if err := mapstructure.Decode(v, ins); err != nil {
		return nil, err
	}
	return ins, nil
}

func DefaultExtractor() *ValueExtractor {
	return &ValueExtractor{
		Type:     _typeMap,
		Child:    make(map[string]*ValueExtractor),
		Callback: make([]func(i interface{}), 0),
	}
}
