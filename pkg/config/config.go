package config

import (
	"gopkg.in/yaml.v2"

	"io/ioutil"
	"log"
)

type Config struct {
	SourceUrl     string   `yaml:"source_url"`
	FileNames     []string `yaml:"filenames"`
	DataPath      string   `yaml:"data_path"`
	Port          int      `yaml:"port"`
	WebFolderPath string   `yaml:"web_folder_path"`
	Db            Db       `yaml:"db"`
}

type Db struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	DatabaseName string `yaml:"database_name"`
	UserName     string `yaml:"user_name"`
	Password     string `yaml:"password"`
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
