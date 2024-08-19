package monitorFolder

import (
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"goConverter/config"
	"goConverter/fileutil"
	"goConverter/hashutil"
	"goConverter/video"
)

// MonitorFolder запускает мониторинг папки для поиска новых файлов
func Run(config *config.Config) {
	currentYear := time.Now().Year()
	yearFolder := filepath.Join(config.SourceVideoFilesFolder, fmt.Sprintf("%d", currentYear))
	daysAgoDuration := time.Duration(-config.DaysAgo) * 24 * time.Hour
	cutoffTime := time.Now().Add(daysAgoDuration)

	if _, err := os.Stat(yearFolder); os.IsNotExist(err) {
		log.Printf("Папка для текущего года не существует: %s", yearFolder)
		return
	}

	err := filepath.Walk(yearFolder, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.ModTime().Before(cutoffTime) {
				return filepath.SkipDir
			}
			return nil
		}
		creationYear := info.ModTime().Year()
		replicaFileFolder := filepath.Join(config.ConvertedFilesFolder, fmt.Sprintf("%d", creationYear), hashutil.Md5Hash(filepath.Dir(path)))

		ProcessFile(path, replicaFileFolder, config)
		return nil
	})

	if err != nil {
		log.Printf("Ошибка обхода папок: %v", err)
	}
}

func ProcessFile(sourceFilePath, convertedFolder string, cfg *config.Config) {
	extension := strings.ToLower(filepath.Ext(sourceFilePath))
	mimeType := mime.TypeByExtension(extension)

	// Извлекаем ключ конвертера из имени файла
	baseFileName := filepath.Base(sourceFilePath)
	fileNameParts := strings.Split(baseFileName, "_")
	converterKey := strings.Split(fileNameParts[len(fileNameParts)-1], ".")[0]

	// Определяем настройки конвертера
	var converterConfig config.ConverterConfig
	if strings.HasPrefix(mimeType, "video/") {
		converterConfig = cfg.Converts["defaultVideo"]
	} else if strings.HasPrefix(mimeType, "image/") {
		converterConfig = cfg.Converts["defaultImage"]
	}

	// Переопределяем настройки, если найден ключ конвертера
	if customConfig, exists := cfg.Converts[converterKey]; exists {
		converterConfig = customConfig
	}

	// Генерируем имя и путь для файла-реплики
	replicaFileName := hashutil.GetReplicaFileName(sourceFilePath, converterConfig.Extension, converterConfig.Resolution)
	replicaFilePath := filepath.Join(convertedFolder, replicaFileName)

	// Проверяем, существует ли файл-реплика
	if _, err := os.Stat(replicaFilePath); err == nil {
		if fileutil.IsFileEmpty(replicaFilePath) {
			log.Printf("Файл существует, но пустой, повторная конвертация: %s", replicaFilePath)
		} else {
			log.Printf("Файл уже существует, пропускаем: %s", replicaFilePath)
			return
		}
	}

	// Создаем папку для файла-реплики
	if err := os.MkdirAll(filepath.Dir(replicaFilePath), 0755); err != nil {
		log.Printf("Ошибка создания папки: %v", err)
		return
	}

	// Выполняем конвертацию в зависимости от типа файла
	var err error
	if strings.HasPrefix(mimeType, "video/") {
		err = video.Convert(sourceFilePath, replicaFilePath, converterConfig)
	} else if strings.HasPrefix(mimeType, "image/") {
		// err = photo.Convert(sourceFilePath, replicaFilePath, converterConfig)
	}

	if err != nil {
		log.Printf("Ошибка конвертации файла %s: %v", sourceFilePath, err)
	}
}
