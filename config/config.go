package config

import (
	"io/ioutil"
	"log"
	"os"
	"path"

	yml "gopkg.in/yaml.v2"
)

type config struct {
	Email string
	Token string
}

// Options ..
var Options config

func (c *config) Read(f string) {
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

func (c *config) Check() {
	if c.Email == "" {
		log.Fatalln("Email not defined")
	}
	if c.Token == "" {
		log.Fatalln("Token not defined")
	}
}
