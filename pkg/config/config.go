package config

import (
	"gopkg.in/yaml.v2"

	"io/ioutil"
	"log"
)

type Config struct {
	Urls                    []string `yaml:"urls"`
	DataPath                string   `yaml:"data_path"`
	DefaultTimeDeltaSeconds string   `yaml:"default_time_delta_seconds"`
	Port                    int      `yaml:"port"`
	OperateLogsPath         string   `yaml:"operate_logs_path"`
	WebFolderPath           string   `yaml:"web_folder_path"`
}

func (c *Config) ReadConfigFromPath(filePath string) *Config {
	yamlFile, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Printf("Error when try ioutil.ReadFile, err: %s", err)
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Printf("Error when try yaml.Unmarshal, err: %s", err)
	}

	return c
}
