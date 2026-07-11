package test

import (
	"fmt"
	hdevConfig "github.com/zhifenghao123/how-dev-go-framework/hdev-config"
	"io/ioutil"
	"reflect"
	"strings"
	"testing"
)

// test config reload
func Test_Reload(t *testing.T) {
	v := 1
	c := func(i interface{}) {
		v = v + i.(int)
	}

	c(2)

	if !reflect.DeepEqual(v, 3) {
		t.Errorf("Test_Reload failed, expect 3, get %v ", v)
	}
}

// test config callback
func Test_Callback(t *testing.T) {
	file := "test-config.yml"
	conf, _ := hdevConfig.Init(hdevConfig.Option{
		File:      file,
		Processor: hdevConfig.GetYamlProcessor,
	})

	var v interface{}

	callback := func(i interface{}) {
		v = i
	}

	// get config value and register callback function
	v = conf.Get("test-config.user.name", callback, callback)
	if !reflect.DeepEqual(v, "hzf") {
		t.Errorf("Test_Callback before failed, expect %v, get %v ", "hzf", v)
	}

	// edit config file
	f, _ := ioutil.ReadFile(file)
	fs := string(f)
	nfs := strings.Replace(fs, "hzf", "yr", 1)
	ioutil.WriteFile(file, []byte(nfs), 0666)
	// rollback after test
	defer ioutil.WriteFile(file, []byte(fs), 0666)

	// reload config file
	conf.Load(file)

	// check callback is effective
	if !reflect.DeepEqual(v, "yr") {
		t.Errorf("Test_Callback after failed, expect %v, get %v ", "yr", v)
	}
}

// test config callback with struct
func Test_CallbackStruct(t *testing.T) {
	file := "test-config.yml"
	conf, _ := hdevConfig.Init(hdevConfig.Option{
		File:      file,
		Processor: hdevConfig.GetYamlProcessor,
	})

	// test struct
	type Struct struct {
		Param1 string
		Param2 string
	}

	var v interface{}

	callback := func(i interface{}) {
		v = i
	}

	// get config value and register callback function
	v = conf.GetAsStruct("test-config.struct", Struct{}, callback)

	fmt.Println(v)

	// edit config file
	f, _ := ioutil.ReadFile(file)
	fs := string(f)
	nfs := strings.Replace(fs, "aaaaaa", "cccccc", 1)
	ioutil.WriteFile(file, []byte(nfs), 0666)
	// rollback after test
	defer ioutil.WriteFile(file, []byte(fs), 0666)

	// reload config file
	conf.Load(file)

	fmt.Println(v)
}

// test config callback with array
func Test_CallbackStructArray(t *testing.T) {
	file := "test-config.yml"
	conf, _ := hdevConfig.Init(hdevConfig.Option{
		File:      file,
		Processor: hdevConfig.GetYamlProcessor,
	})

	// test struct
	type Struct struct {
		Param1 string
		Param2 string
	}

	var v interface{}

	callback := func(i interface{}) {
		v = i
	}

	// get config value and register callback function
	v = conf.GetAsStructArray("test-config.struct_array", Struct{}, callback)

	for _, item := range v.([]interface{}) {
		fmt.Println(item)
	}

	//edit config file
	f, _ := ioutil.ReadFile(file)
	fs := string(f)
	nfs := strings.Replace(fs, "aaaaaaaa", "aaaaaaaa1", 1)
	ioutil.WriteFile(file, []byte(nfs), 0666)
	// rollback after test
	defer ioutil.WriteFile(file, []byte(fs), 0666)

	// reload config file
	conf.Load(file)

	for _, item := range v.([]interface{}) {
		fmt.Println(item)
	}
}
