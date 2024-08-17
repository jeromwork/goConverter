package hashutil

import (
	"crypto/md5"
	"encoding/hex"
)

// Md5Hash получает хэш MD5 строки
func Md5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

// GetReplicaFileName получает имя файла для сохранения конвертированного видео
func GetReplicaFileName(sourceVideoPath, extension, resolution string) string {
	hash := Md5Hash(sourceVideoPath + extension + resolution)
	if len(hash) > 12 {
		hash = hash[len(hash)-12:] // Оставляем только последние 12 символов
	}
	return hash + "." + extension
}
