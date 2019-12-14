package conf

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Port string `yaml:"port"`
	Log  struct {
		Level  string `yaml:"level"`
		Output string `yaml:"output"`
	}
	Database struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"dbname"`
		Host     string `yaml:"host"`
	}
	CodeBaseDir       string `yaml:"codeBaseDir"`
	SubmissionBaseDir string `yaml:"submissionBaseDir"`
	Redis             struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
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
