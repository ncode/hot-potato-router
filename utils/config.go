package utils

import (
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"strings"
)

type Config struct {
	Options map[string]map[string]string
}

func NewConfig() Config {
	cfg := Config{}

	config_file := os.Getenv("HPR_CONF")
	if strings.EqualFold(config_file, "") {
		config_file = "/etc/hpr/hpr.yml"
	}
	f, err := ioutil.ReadFile(config_file)
	CheckPanic(err, "Unable to open YAML file")

	err = goyaml.Unmarshal(f, &cfg.Options)
	CheckPanic(err, "Unable to parse YAML file")
	return cfg
}
