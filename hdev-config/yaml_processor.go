package hdevconfig

import (
	"gopkg.in/yaml.v2"
	"strconv"
)

var _ IProcessor = YamlProcessor{}

// YamlProcessor processor type of yaml
type YamlProcessor struct {
	FileReader
}

// Load readFile a yaml config file to map data
func (y YamlProcessor) LoadFileContentAsMapData(file string) (map[string]interface{}, error) {
	dat, err := y.readFile(file)
	if err != nil {
		return nil, err
	}
	var d map[string]interface{}
	err = yaml.Unmarshal(dat, &d)
	if err != nil {
		return nil, err
	}

	return yamlTrans(d), nil
}

// transform the interface{} of map[interface{}] interface{} to map[string]interface{}
func yamlTrans(m map[string]interface{}) map[string]interface{} {
	for k, v := range m {
		if v1, ok := v.(map[interface{}]interface{}); ok {
			mm := make(map[string]interface{})
			for k1, v2 := range v1 {
				if vk1, ok := k1.(string); ok {
					mm[vk1] = v2
					continue
				}

				if vk2, ok := k1.(int); ok {
					mm[strconv.Itoa(vk2)] = v2
					continue
				}
			}
			m[k] = yamlTrans(mm)
		}
	}
	return m
}

// YamlProcessor returns a Processor which readFile yaml config
func GetYamlProcessor() IProcessor {
	return &YamlProcessor{}
}
