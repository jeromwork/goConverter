package video

import (
	"fmt"
	"goConverter/config"
	"goConverter/fileutil"
	"log"
	"os/exec"
	"path/filepath"
)

// Convert выполняет конвертацию видео с учетом переданных параметров
func Convert(sourcePath, targetPath string, converterConfig config.ConverterConfig) error {
	// Получаем параметры из конфигурации
	resolution := converterConfig.Resolution
	extension := converterConfig.Extension
	watermarkImage := converterConfig.WatermarkImage
	watermarkPosition := converterConfig.WatermarkPosition

	if watermarkImage == "" && converterConfig.Watermark != "" {
		watermarkImage = "watermarkEastclinicWhite.png"
	}
	if watermarkPosition == "" && (watermarkImage != "" || converterConfig.Watermark != "") {
		watermarkPosition = "W-w-20:H-h-20"
	}

	// Проверка разрешения и получение параметра масштабирования
	scale := fileutil.ResolutionToScale(resolution)
	if scale == "" {
		return fmt.Errorf("неизвестное разрешение: %s", resolution)
	}

	// Изменение расширения целевого файла
	if filepath.Ext(targetPath) != "."+extension {
		targetPath = targetPath[:len(targetPath)-len(filepath.Ext(targetPath))] + "." + extension
	}

	// Получение абсолютного пути к водяному знаку
	absWatermarkPath, err := filepath.Abs(watermarkImage)
	if err != nil {
		return fmt.Errorf("не удалось получить абсолютный путь к изображению водяного знака: %v", err)
	}

	// Основной фильтр масштабирования
	filterComplex := fmt.Sprintf(
		"[0:v]scale=%s,setsar=1[scaled];[1:v]scale=w=-1:h=ih*0.3[wm]",
		scale,
	)

	// Добавляем фильтр наложения водяного знака, если требуется
	if watermarkPosition != "" {
		filterComplex += fmt.Sprintf(";[scaled][wm]overlay=%s", watermarkPosition)
	} else {
		filterComplex += ";[scaled][wm]overlay=shortest=1"
	}

	// Логирование фильтра для отладки
	log.Printf("Используемый фильтр complex: %s\n", filterComplex)

	cmdArgs := []string{
		"-y", "-i", sourcePath,
		"-i", absWatermarkPath,
		"-filter_complex", filterComplex,
		targetPath,
	}

	// Выполнение команды ffmpeg
	cmd := exec.Command("ffmpeg", cmdArgs...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("Ошибка конвертации %s: %v\n", sourcePath, err)
		log.Printf("Вывод ffmpeg: %s\n", string(output))
		return err
	}

	// Проверка, что файл не пустой
	if fileutil.IsFileEmpty(targetPath) {
		log.Printf("Файл %s пустой после конвертации, повторная конвертация...\n", targetPath)
		return fmt.Errorf("пустой файл после конвертации")
	}

	log.Printf("Успешная конвертация %s в %s с разрешением %s\n", sourcePath, targetPath, resolution)
	return nil
}
