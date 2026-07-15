package hdevconfig

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"sync"
)

type Option struct {
	File      string
	Processor ProcessorConstruct
	Extractor ExtractorConstruct
}

type Config struct {
	processor IProcessor
	extractor IExtractor
	watcher   *fsnotify.Watcher
	l         sync.RWMutex
}

func Init(option Option) (c *Config, err error) {
	c = &Config{
		extractor: DefaultExtractor(),
	}

	if option.Processor != nil {
		c.processor = option.Processor()
	}

	if option.Extractor != nil {
		c.extractor = option.Extractor()
	}

	if option.File != "" {
		err = c.Load(option.File)
	}

	return c, err
}

func (c *Config) Get(key string, b ...func(i interface{})) interface{} {
	return c.extractor.Get(key, b...)
}

func (c *Config) GetAsInt(key string, b ...func(i interface{})) int {
	return c.extractor.GetAsInt(key, b...)
}

func (c *Config) GetAsFloat(key string, b ...func(i interface{})) float64 {
	return c.extractor.GetAsFloat(key, b...)
}

func (c *Config) GetAsString(key string, b ...func(i interface{})) string {
	return c.extractor.GetAsString(key, b...)
}

func (c *Config) GetAsBool(key string, b ...func(i interface{})) bool {
	return c.extractor.GetAsBool(key, b...)
}

func (c *Config) GetAsArray(key string, b ...func(i interface{})) []interface{} {
	return c.extractor.GetAsArray(key, b...)
}

func (c *Config) GetAsMap(key string, b ...func(i interface{})) map[string]interface{} {
	return c.extractor.GetAsMap(key, b...)
}

func (c *Config) GetAsStruct(key string, s interface{}, b ...func(i interface{})) interface{} {
	return c.extractor.GetAsStruct(key, s, b...)
}

func (c *Config) GetAsStructArray(key string, s interface{}, b ...func(i interface{})) []interface{} {
	return c.extractor.GetAsStructArray(key, s, b...)
}

func (c *Config) Load(file string, alias ...string) (e error) {
	if e = c.watch(file); e != nil {
		return
	}
	c.l.Lock()
	defer c.l.Unlock()
	if len(alias) > 0 {
		fileAlias[file] = alias[0]
	}
	return c.reload(file)
}

func (c *Config) reload(file string) (e error) {
	return Loader(func() (IProcessor, IExtractor) {
		return c.processor, c.extractor
	}).load(file)
}

func (c *Config) watch(file string) (e error) {
	if e = c.checkWatcher(); e != nil {
		return
	}

	return c.watcher.Add(file)
}

func (c *Config) checkWatcher() (e error) {
	if c.watcher != nil {
		return
	}

	if c.watcher, e = fsnotify.NewWatcher(); e != nil {
		return
	}

	go func() {
		for {
			select {
			case event := <-c.watcher.Events:
				log.Println("config file modified : ", event.Name)
				if event.Op&fsnotify.Write == fsnotify.Write {
					err := c.reload(event.Name)
					if err != nil {
						log.Println("config file reload error : ", err)
					}
				}
			case err := <-c.watcher.Errors:
				log.Println("error : ", err)
			}
		}
	}()

	return
}
