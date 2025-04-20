package Error

import (
	conf "back/config"
	"back/internal/utils"
	"encoding/json"
	"io"
	"log"
	"os"
)

type Pars struct {
	Error  string  `json:"error"`
	TimErr string  `json:"time"`
}

func GetErr(err error) bool{
	if err == nil {
		return false
	}

	log.Println(err)

	pars := Pars{
		Error:  err.Error(),
		TimErr: utils.TimeFormat(),
	}

	file, err := utils.FileEx(conf.Cfg.ErrorPath)
	if err != nil {
		file, err = os.Create(conf.Cfg.ErrorPath)
		if err != nil {
			log.Println("Ошибка открытия файла:", err)
			return true
		}
	}
	defer file.Close()

	var logs []Pars
	data, _ := io.ReadAll(file)
	if len(data) > 0 {
		if unmarshalErr := json.Unmarshal(data, &logs); unmarshalErr != nil {
			log.Println("Ошибка парсинга JSON:", unmarshalErr)
		}
	}

	logs = append([]Pars{pars}, logs...)

	if _, seekErr := file.Seek(0, 0); seekErr != nil {
		log.Println("Ошибка seek файла:", seekErr)
		return true
	}

	if truncateErr := file.Truncate(0); truncateErr != nil {
		log.Println("Ошибка очистки файла:", truncateErr)
		return true
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if encodeErr := encoder.Encode(logs); encodeErr != nil {
		log.Println("Ошибка записи в файл:", encodeErr)
	} else {
		log.Println("Данные успешно записаны в error_log.json")
	}
	return false
}
