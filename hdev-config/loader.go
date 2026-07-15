package hdevconfig

import (
	"path"
	"path/filepath"
	"strings"
)

var fileAlias map[string]string

func init() {
	fileAlias = make(map[string]string)
}

// ILoader interface of loader
type ILoader interface {
	load(string) error
}

// Loader loader
type Loader func() (IProcessor, IExtractor)

// load config with referred processor and put into extractor
func (l Loader) load(file string) error {
	p, e := l()
	if p == nil {
		p = getProcessor(file)
	}
	data, err := p.LoadFileContentAsMapData(file)
	if err != nil {
		return err
	}

	if name, ok := fileAlias[file]; ok {
		e.Load(data, name)
		return nil
	}

	e.Load(data, Filename(file))
	return nil
}

// Filename get filename
func Filename(file string) string {
	filename := path.Base(filepath.ToSlash(file))
	suffix := path.Ext(filename)
	return strings.TrimSuffix(filename, suffix)
}

// get processor by config file extension
func getProcessor(file string) IProcessor {
	ext := path.Ext(file)
	switch strings.ToLower(ext) {
	case _yaml, _yml:
		return GetYamlProcessor()
	case _json:
		return GetJsonProcessor()
	case _ini:
		return GetIniProcessor()
	default:
		return GetJsonProcessor()
	}
}
