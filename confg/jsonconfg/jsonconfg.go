package jsonconfg

import (
	"encoding/json"
	"github.com/kidstuff/toys/confg"
	"os"
)

func init() {
	confg.Register("jsonconfg", &JSONConfig{})
}

type JSONConfig struct {
	data map[string]interface{}
	file *os.File
}

func (c *JSONConfig) Load(path string) (err error) {
	c.file, err = os.Open(path)
	if err != nil {
		return
	}
	c.data = make(map[string]interface{})
	dec := json.NewDecoder(c.file)
	err = dec.Decode(&c.data)
	return
}

func (c *JSONConfig) Close() error {
	return c.file.Close()
}

func (c *JSONConfig) Set(k string, v interface{}) {
	c.data[k] = v
	//TODO: make change to file
}

func (c *JSONConfig) Get(k string) interface{} {
	return c.data[k]
}

func (c *JSONConfig) Del(k string) {
	delete(c.data, k)
	//TODO: make change to file
}

var _ confg.Configurator = &JSONConfig{}
