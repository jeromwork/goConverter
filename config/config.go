package config

import (
	"encoding/json"
	"io/ioutil"
)

// Config структура для хранения настроек
type Config struct {
	SourceVideoFilesFolder string   `json:"sourceVideoFilesFolder"`
	ConvertedFilesFolder   string   `json:"convertedFilesFolder"`
	DaysAgo                int      `json:"daysAgo"`
	Formats                []Format `json:"formats"`
}

// Format структура для хранения форматов и разрешений
type Format struct {
	Extension   string   `json:"extension"`
	Resolutions []string `json:"resolutions"`
}

// LoadConfig загружает настройки из файла
func LoadConfig(configPath string) (*Config, error) {
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	err = json.Unmarshal(file, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
