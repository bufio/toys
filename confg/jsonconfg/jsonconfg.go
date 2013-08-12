package jsonconfg

import (
	"encoding/json"
	"github.com/kidstuff/toys/confg"
	"github.com/kidstuff/toys/util/errs"
	"os"
)

func init() {
	confg.Register("jsonconfg", &JSONConfig{})
}

type JSONConfig struct {
	data map[string]interface{}
	file *os.File
}

func (c *JSONConfig) Load(path string) error {
	var err error
	c.file, err = os.Open(path)
	if err != nil {
		return errs.Errf(err, "jsonconfg: cannot load the file: %s", path)
	}

	c.data = make(map[string]interface{})
	dec := json.NewDecoder(c.file)
	err = dec.Decode(&c.data)
	if err != nil {
		return errs.Errf(err, "jsonconfig: cannot deocde data in %s", path)
	}

	return nil
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
