package config

import (
	"encoding/json"
	"io/ioutil"
)

// Config структура для хранения настроек
type Config struct {
	SourceVideoFilesFolder string                     `json:"sourceVideoFilesFolder"`
	ConvertedFilesFolder   string                     `json:"convertedFilesFolder"`
	DaysAgo                int                        `json:"daysAgo"`
	Converts               map[string]ConverterConfig `json:"converts"`
}

// ConverterConfig структура для хранения настроек конвертера
type ConverterConfig struct {
	Extension         string `json:"extension"`
	Resolution        string `json:"resolution"`
	Watermark         string `json:"watermark"`
	WatermarkImage    string `json:"watermarkImage"`
	WatermarkPosition string `json:"watermarkPosition"` // default, W-w-20:H-h-20
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
