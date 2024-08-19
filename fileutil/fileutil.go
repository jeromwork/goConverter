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

// resolutionToScale преобразует разрешение в параметры масштабирования для ffmpeg
func ResolutionToScale(resolution string) string {
	switch resolution {
	case "360p":
		return "w=-1:h=360"
	case "720p":
		return "w=-1:h=720"
	case "1080p":
		return "w=-1:h=1080"
	default:
		return "" // Для неизвестного разрешения вернем пустую строку
	}
}
