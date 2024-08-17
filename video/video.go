package video

import (
	"fmt"
	"log"
	"os/exec"

	"goConverter/fileutil"
)

// ConvertVideo конвертирует видео с использованием ffmpeg
func Convert(sourcePath, targetPath, resolution string) error {
	scale := resolutionToScale(resolution)
	if scale == "" {
		return fmt.Errorf("неизвестное разрешение: %s", resolution)
	}

	cmd := exec.Command("ffmpeg", "-y", "-i", sourcePath, "-vf", "scale="+scale, targetPath)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Ошибка конвертации %s: %v\n", sourcePath, err)
		log.Printf("Вывод ffmpeg: %s\n", string(output))
		return err
	}

	if fileutil.IsFileEmpty(targetPath) {
		log.Printf("Файл %s пустой после конвертации, повторная конвертация...\n", targetPath)
		return fmt.Errorf("пустой файл после конвертации")
	}

	log.Printf("Успешная конвертация %s в %s с разрешением %s\n", sourcePath, targetPath, resolution)
	return nil
}

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
