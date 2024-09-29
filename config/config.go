package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Config struct {
	WatchFolder       string `yaml:"watch_folder"`
	IncludeSubfolders bool   `yaml:"include_subfolders"`
	MetadataFile      string `yaml:"metadata_file"`
	SyncDelaySeconds  int    `yaml:"sync_delay_seconds"` // TODO
	LogLevel          string `yaml:"log_level"`          // TODO
}

func LoadConfig(configPath string) (*Config, error) {
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
