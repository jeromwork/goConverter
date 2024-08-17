package fileutil

import (
	"log"
	"os"
)

// IsFileEmpty проверяет размер файла
func IsFileEmpty(filePath string) bool {
	info, err := os.Stat(filePath)
	if err != nil {
		log.Printf("Ошибка проверки файла %s: %v", filePath, err)
		return true
	}
	return info.Size() == 0
}
