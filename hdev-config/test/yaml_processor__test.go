package test

import (
	"fmt"
	hdevConfig "github.com/zhifenghao123/how-dev-go-framework/hdev-config"
	"testing"
)

// test config load of yaml file
func Test_Yaml_Load(t *testing.T) {
	yaml := hdevConfig.GetYamlProcessor()
	data, err := yaml.LoadFileContentAsMapData("test-config.yml")
	fmt.Println(data)
	if err != nil {
		t.Errorf("Test_Yaml_Load failed, %s", err.Error())
	}
}
