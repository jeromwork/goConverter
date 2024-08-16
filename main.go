package main

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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

// Загружаем настройки из файла
func loadConfig(configPath string) (*Config, error) {
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

// Получение хэша MD5 строки
func md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// Получение имени файла для сохранения конвертированного видео
func getReplicaFileName(sourceVideoPath, extension, resolution string) string {
	hash := md5Hash(sourceVideoPath + extension + resolution)
	return fmt.Sprintf("%s.%s", hash, extension)
}

// Преобразование разрешения в формат width:height
func resolutionToScale(resolution string) string {
	switch resolution {
	case "360p":
		return "640:360"
	case "720p":
		return "1280:720"
	case "1080p":
		return "1920:1080"
	default:
		return "" // Для неизвестного разрешения вернем пустую строку
	}
}

// Проверка размера файла
func isFileEmpty(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Ошибка проверки файла %s: %v", filePath, err)
		return true
	}
	return info.Size() == 0
}

// Конвертация видео с использованием ffmpeg
func convertVideo(sourcePath, targetPath, resolution string) error {
	scale := resolutionToScale(resolution)
	if scale == "" {
		return fmt.Errorf("неизвестное разрешение: %s", resolution)
	}

	// Пример команды для конвертации видео с помощью ffmpeg с флагом -y для перезаписи файла
	cmd := exec.Command("ffmpeg", "-y", "-i", sourcePath, "-vf", "scale="+scale, targetPath)

	// Запускаем команду и ждем завершения
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Ошибка конвертации %s: %v\n", sourcePath, err)
		log.Printf("Вывод ffmpeg: %s\n", string(output))
		return err
	}

	if isFileEmpty(targetPath) {
		log.Printf("Файл %s пустой после конвертации, повторная конвертация...\n", targetPath)
		return fmt.Errorf("пустой файл после конвертации")
	}

	log.Printf("Успешная конвертация %s в %s с разрешением %s\n", sourcePath, targetPath, resolution)
	return nil
}

// Обработка видеофайла
func processVideoFile(sourceFilePath, convertedFolder, parentFolder string, formats []Format) {
	fileName := strings.TrimSuffix(filepath.Base(sourceFilePath), filepath.Ext(sourceFilePath))
	for _, format := range formats {
		for _, resolution := range format.Resolutions {
			replicaFileName := getReplicaFileName(sourceFilePath, format.Extension, resolution)
			replicaFilePath := filepath.Join(convertedFolder, parentFolder, fileName, replicaFileName)

			if _, err := os.Stat(replicaFilePath); err == nil {
				if isFileEmpty(replicaFilePath) {
					log.Printf("Файл существует, но пустой, повторная конвертация: %s", replicaFilePath)
				} else {
					log.Printf("Файл уже существует, пропускаем: %s", replicaFilePath)
					continue
				}
			}

			if err := os.MkdirAll(filepath.Dir(replicaFilePath), 0755); err != nil {
				log.Printf("Ошибка создания папки: %v", err)
				continue
			}

			for {
				err := convertVideo(sourceFilePath, replicaFilePath, resolution)
				if err == nil {
					break // Если конвертация прошла успешно, выходим из цикла
				}
				log.Printf("Попытка повторной конвертации %s\n", sourceFilePath)
			}
		}
	}
}

// Мониторинг папки
func monitorFolder(config *Config) {
	sourceDir := config.SourceVideoFilesFolder
	daysAgoDuration := time.Duration(-config.DaysAgo) * 24 * time.Hour
	cutoffTime := time.Now().Add(daysAgoDuration)

	err := filepath.Walk(sourceDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.ModTime().Before(cutoffTime) {
				return filepath.SkipDir
			}
			return nil
		}

		parentFolder := md5Hash(filepath.Dir(path))
		processVideoFile(path, config.ConvertedFilesFolder, parentFolder, config.Formats)
		return nil
	})

	if err != nil {
		log.Printf("Ошибка обхода папок: %v", err)
	}
}

func main() {
	configPath := flag.String("config", "settings.json", "Путь к файлу конфигурации")
	flag.Parse()

	config, err := loadConfig(*configPath)
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	log.Println("Запуск мониторинга...")
	for {
		monitorFolder(config)
		time.Sleep(10 * time.Second) // Интервал проверки новых файлов
	}
}
