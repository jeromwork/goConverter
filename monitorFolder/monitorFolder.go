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

		ProcessFile(path, replicaFileFolder, config.Formats)
		return nil
	})

	if err != nil {
		log.Printf("Ошибка обхода папок: %v", err)
	}
}

func ProcessFile(sourceFilePath, convertedFolder string, formats []config.Format) {
	for _, format := range formats {
		for _, resolution := range format.Resolutions {
			replicaFileName := hashutil.GetReplicaFileName(sourceFilePath, format.Extension, resolution)
			replicaFilePath := filepath.Join(convertedFolder, replicaFileName)

			if _, err := os.Stat(replicaFilePath); err == nil {
				if fileutil.IsFileEmpty(replicaFilePath) {
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

			extension := strings.ToLower(filepath.Ext(sourceFilePath))
			mimeType := mime.TypeByExtension(extension)
			var err error
			if strings.HasPrefix(mimeType, "video/") {
				err = video.Convert(sourceFilePath, replicaFilePath, resolution)
			} else if strings.HasPrefix(mimeType, "image/") {
				err = video.Convert(sourceFilePath, replicaFilePath, resolution)
			} else {
				log.Printf("Неподдерживаемый тип файла: %s", sourceFilePath)
			}

			if err == nil {
				break // Если конвертация прошла успешно, выходим из цикла
			}
			// for {
			// 	err := ConvertVideo(sourceFilePath, replicaFilePath, resolution)
			// 	if err == nil {
			// 		break // Если конвертация прошла успешно, выходим из цикла
			// 	}
			// 	log.Printf("Попытка повторной конвертации %s\n", sourceFilePath)
			// }
		}
	}
}
