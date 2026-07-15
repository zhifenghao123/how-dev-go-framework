package hdevconfig

import (
	"encoding/json"
	"github.com/json-iterator/go"
)

var j jsoniter.API

func init() {
	j = jsoniter.Config{
		EscapeHTML:             false,
		SortMapKeys:            false,
		ValidateJsonRawMessage: true,
		UseNumber:              true,
	}.Froze()
}

var _ IProcessor = JsonProcessor{}

// JsonProcessor processor type of json
type JsonProcessor struct {
	FileReader
}

// LoadFileContentAsMapData 读取json配置文件为map数据
func (y JsonProcessor) LoadFileContentAsMapData(file string) (map[string]interface{}, error) {
	dat, err := y.readFile(file)
	if err != nil {
		return nil, err
	}
	var d map[string]interface{}
	if err = j.Unmarshal(dat, &d); err != nil {
		return nil, err
	}
	d = jsonTrans(d)
	return d, nil
}

// transfer json data
func jsonTrans(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		if v1, ok := v.(map[string]interface{}); ok {
			m[k] = jsonTrans(v1)
			continue
		}
		if v1, ok := v.([]interface{}); ok {
			m[k] = jsonArray(v1)
			continue
		}
		if v1, ok := v.(json.Number); ok {
			m[k] = jsonNumber(v1)
		}
	}
	return m
}

// transfer json array
func jsonArray(a []interface{}) interface{} {
	if len(a) == 0 {
		return a
	}

	if _, ok := a[0].(json.Number); ok {
		var arr []interface{}
		for _, v1 := range a {
			arr = append(arr, jsonNumber(v1.(json.Number)))
		}
		return arr
	}

	if _, ok := a[0].(map[string]interface{}); ok {
		var arr []map[string]interface{}
		for _, v1 := range a {
			arr = append(arr, jsonTrans(v1.(map[string]interface{})))
		}
		return arr
	}

	return a
}

// json number to int
func jsonNumber(n json.Number) interface{} {
	i, err := n.Int64()
	if err == nil {
		return int(i)
	}
	f, err := n.Float64()
	if err == nil {
		return f
	}
	return n.String()
}

// GetJsonProcessor returns a Processor which readFile json config
func GetJsonProcessor() IProcessor {
	return &JsonProcessor{}
}
