package config

import (
	"bytes"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Commands map[string][]string `yaml:"commands"`
}

var Conf Config

func Load() {
	f, err := os.Open("/opt/config/talenesia.yaml")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var buf bytes.Buffer
	_, err = io.Copy(&buf, f)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(buf.Bytes(), &Conf)
	if err != nil {
		panic(err)
	}
}
