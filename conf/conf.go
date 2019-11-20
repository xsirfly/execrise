package conf

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Config struct {
	Port string `yaml:"port"`
	Log  struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	}
	Database struct {
		User string `yaml:"user"`
		Password string `yaml:"password"`
		DBName string `yaml:"dbname"`
		Host string `yaml:"host"`
	}
}

var config Config

func Init(env string) error {
	configFile := "conf_" + env + ".yml"
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(b, &config); err != nil {
		return err
	}
	config.Port = ":" + config.Port
	return nil
}

func GetConf() *Config {
	return &config
}
