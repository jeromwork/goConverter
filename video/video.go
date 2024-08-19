package video

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"goConverter/config"
	"goConverter/fileutil"
)

// Convert выполняет конвертацию видео или изображения с учетом переданных параметров
func Convert(sourcePath, targetPath string, converterConfig config.ConverterConfig) error {
	// Получаем параметры из конфигурации
	resolution := converterConfig.Resolution
	extension := converterConfig.Extension
	watermarkImage := converterConfig.WatermarkImage
	// Дополнительные параметры для ffmpeg можно добавить здесь
	// Пример: битрейт, частота кадров и т.д.
	// additionalParams := converterConfig["additionalParams"].(string) // например, "-b:v 1M"

	scale := resolutionToScale(resolution)
	if scale == "" {
		return fmt.Errorf("неизвестное разрешение: %s", resolution)
	}

	// Изменяем расширение целевого файла на указанное в конфигурации (если необходимо)
	if filepath.Ext(targetPath) != "."+extension {
		targetPath = targetPath[:len(targetPath)-len(filepath.Ext(targetPath))] + "." + extension
	}

	if watermarkImage == "" {
		watermarkImage = "watermarkEastclinicWhite.png"
	}
	absWatermarkPath, err := filepath.Abs(watermarkImage)
	if err != nil {
		return fmt.Errorf("не удалось получить абсолютный путь к изображению водяного знака: %v", err)
	}

	// Команду для ffmpeg строим с учетом дополнительных параметров
	cmdArgs := []string{
		"-y", "-i", sourcePath,
		"-i", absWatermarkPath,
		"-filter_complex", "scale=" + scale + ",overlay=W-w-30:H-h-30",
		targetPath,
	}

	// Команда для ffmpeg
	cmd := exec.Command("ffmpeg", cmdArgs...)

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

// resolutionToScale преобразует разрешение в параметры масштабирования для ffmpeg
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
