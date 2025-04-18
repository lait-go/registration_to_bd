package utils

import (
	"log"
	"os"
	"time"
)

func FileEx(filename string) (*os.File, error) {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, err
		}
	}

	return file, nil
}

func FileClear(file *os.File) bool {
	if _, seekErr := file.Seek(0, 0); seekErr != nil {
		log.Println("Ошибка seek файла:", seekErr)
		return true
	}

	if truncateErr := file.Truncate(0); truncateErr != nil {
		log.Println("Ошибка очистки файла:", truncateErr)
		return true
	}
	return false
}

func TimeFormat() string {
	return time.Now().Format("2006-01-02 15:04:05")
}