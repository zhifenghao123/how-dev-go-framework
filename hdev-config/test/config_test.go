package test

import (
	hdevConfig "github.com/zhifenghao123/how-dev-go-framework/hdev-config"
	"reflect"
	"testing"
)

// test config of yaml
func Test_Yaml(t *testing.T) {
	conf, _ := hdevConfig.Init(hdevConfig.Option{
		File: "test-config.yml",
	})

	//test struct
	type Struct struct {
		Param1 string
		Param2 string
	}

	sa := conf.GetAsStructArray("test-config.struct_array", Struct{})

	tests := []struct {
		want   interface{}
		result interface{}
	}{
		{nil, conf.Get("")},
		{"hzf", conf.Get("test-config.user.name")},
		{nil, conf.Get("test-config.user.name1")},
		{"hzf", conf.GetAsString("test-config.user.name")},
		{"", conf.GetAsString("test-config.user.name1")},
		{"30", conf.GetAsString("test-config.user.age")},
		{30, conf.GetAsInt("test-config.user.age")},
		{0, conf.GetAsInt("test-config.user.age1")},
		{true, conf.GetAsBool("test-config.user.male")},
		{false, conf.GetAsBool("test-config.user.female")},
		{false, conf.GetAsBool("test-config.user.female1")},
		{[]interface{}{"aaa", "bbb", "ccc"}, conf.GetAsArray("test-config.user.list")},
		{map[string]interface{}{"aa": "aa", "bb": "bb", "11": 11}, conf.GetAsMap("test-config.user.property")},
		{&Struct{Param1: "aaaaaa", Param2: "bbbbbb"}, conf.GetAsStruct("test-config.struct", Struct{})},
		{&Struct{Param1: "aaaaaaaa", Param2: "bbbbbbbb"}, sa[0]},
		{&Struct{Param1: "cccccccc", Param2: "dddddddd"}, sa[1]},
	}

	for i, test := range tests {
		if !reflect.DeepEqual(test.want, test.result) {
			t.Errorf("Test_Yaml %d failed, expect %v, get %v ", i, test.want, test.result)
		}
	}

	if v := conf.GetAsArray("test-config.user.list1"); v != nil {
		t.Errorf("Test_Yaml conf.GetAsGetAsArray failed, expect %v, get %v ", nil, v)
	}

	if v := conf.GetAsMap("test-config.user.property1"); v != nil {
		t.Errorf("Test_Yaml conf.GetAsGetAsMap failed, expect %v, get %v ", nil, v)
	}

	if v := conf.GetAsStruct("test-config.struct1", Struct{}); v != nil {
		t.Errorf("Test_Yaml conf.GetAsGetAsStruct failed, expect %v, get %v ", nil, v)
	}

	if v := conf.GetAsStructArray("test-config.struct_array1", Struct{}); v != nil {
		t.Errorf("Test_Yaml conf.GetAsGetAsStructArray failed, expect %v, get %v ", nil, v)
	}

}

// test config of json
func Test_Json(t *testing.T) {
	conf, _ := hdevConfig.Init(hdevConfig.Option{
		File: "test-config.json",
	})

	//test struct
	type Struct struct {
		Param1 string
		Param2 string
		Param3 float64
	}

	sa := conf.GetAsStructArray("test-config.struct_array", Struct{})

	tests := []struct {
		want   interface{}
		result interface{}
	}{
		{"hzf", conf.Get("test-config.user.name")},
		{"hzf", conf.GetAsString("test-config.user.name")},
		{30, conf.GetAsInt("test-config.user.age")},
		{true, conf.GetAsBool("test-config.user.male")},
		{false, conf.GetAsBool("test-config.user.female")},
		{[]interface{}{"aaa", "bbb", "ccc"}, conf.GetAsArray("test-config.user.list")},
		{map[string]interface{}{"aa": "aa", "bb": "bb", "11": 11}, conf.GetAsMap("test-config.user.property")},
		{&Struct{Param1: "aaa", Param2: "bbb", Param3: 0}, conf.GetAsStruct("test-config.struct", Struct{})},
		{&Struct{Param1: "aaa", Param2: "bbb", Param3: 10}, sa[0]},
		{&Struct{Param1: "ccc", Param2: "ddd", Param3: 3.3}, sa[1]},
	}

	for i, test := range tests {
		if !reflect.DeepEqual(test.want, test.result) {
			t.Errorf("Test_Yaml %d failed, expect %v, get %v ", i, test.want, test.result)
		}
	}
}

// test config of ini
func Test_Ini(t *testing.T) {
	conf, _ := hdevConfig.Init(hdevConfig.Option{
		File: "test-config.ini",
	})

	//test struct
	type Struct struct {
		Param1 string
		Param2 string
	}

	//sa := conf.GetAsGetAsStructArray("test-config.struct_array", GetAsStruct{})

	tests := []struct {
		want   interface{}
		result interface{}
	}{
		{"hzf", conf.Get("test-config.user.name")},
		{"aa", conf.GetAsString("test-config.user.property.aa")},
		{30, conf.GetAsInt("test-config.user.age")},
		{11, conf.GetAsInt("test-config.user.property.11")},
		{true, conf.GetAsBool("test-config.user.male")},
		{false, conf.GetAsBool("test-config.user.female")},
		{map[string]interface{}{"aa": "aa", "bb": "bb", "11": "11"}, conf.GetAsMap("test-config.user.property")},
	}

	for i, test := range tests {
		if !reflect.DeepEqual(test.want, test.result) {
			t.Errorf("Test_Yaml %d failed, expect %v, get %v ", i, test.want, test.result)
		}
	}

}
