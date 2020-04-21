package spl

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	yml "gopkg.in/yaml.v2"
)

// Config ..
type Config struct {
	Email    string
	Token    string
	BaseURL  string
	JWTtoken string
}

// Conf ..
var Conf Config

func (c *Config) Read(f string) {
	if !path.IsAbs(f) && f[:1] == "~" {
		f = path.Join(os.Getenv("HOME"), f[1:])
	}

	bytes, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatalf("Read file %s is failed: %s", f, err)
	}

	err = yml.Unmarshal(bytes, c)
	if err != nil {
		log.Fatalf("Read config is failed: %s", err)
	}
}
