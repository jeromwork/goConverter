package image

import (
	"fmt"
	"goConverter/config"
	"goConverter/fileutil"
	"os/exec"
	"path/filepath"
)

func Convert(sourcePath, targetPath string, converterConfig config.ConverterConfig) error {
	watermark := converterConfig.Watermark
	watermarkImage := converterConfig.WatermarkImage
	watermarkPosition := converterConfig.WatermarkPosition // Добавляем позицию водяного знака
	if watermarkPosition == "" && (watermarkImage != "" || watermark != "") {
		watermarkPosition = "W-w-20:H-h-20"
	}
	cmdArgs := []string{"-y", "-i", sourcePath}
	resolution := converterConfig.Resolution
	if watermarkImage == "" {
		watermarkImage = "watermarkEastclinicWhite.png"
	}
	absWatermarkPath, err := filepath.Abs(watermarkImage)
	if err != nil {
		return err
	}
	if absWatermarkPath != "" {
		cmdArgs = append(cmdArgs, "-i", absWatermarkPath)
	}

	scale := fileutil.ResolutionToScale(resolution)
	if scale == "" {
		return fmt.Errorf("неизвестное разрешение: %s", resolution)
	}

	// Добавляем фильтр масштабирования
	filterComplex := fmt.Sprintf(
		"[0:v]scale=%s,setsar=1[scaled];[1:v]scale=w=-1:h=ih*0.3[wm]",
		scale,
	)
	if watermarkImage != "" {
		// Добавляем фильтр наложения водяного знака
		filterComplex += fmt.Sprintf(";[scaled][wm]overlay=%s", watermarkPosition)
	}

	cmdArgs = append(cmdArgs, "-filter_complex", filterComplex, targetPath)

	// Создаем команду для наложения водяного знака и изменения размера изображения
	cmd := exec.Command("ffmpeg", cmdArgs...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ошибка конвертации изображения: %v, вывод: %s", err, string(output))
	}

	return nil
}
