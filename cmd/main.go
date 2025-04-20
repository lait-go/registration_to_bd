package main

import (
	conf "back/config"
	"back/db"
	handler "back/internal/api"
	Error "back/internal/err"
	LogWork "back/internal/utils/log"
	"fmt"
	"net/http"
)

func main() {
	conf.Cfg = conf.Config_work()

	logger := LogWork.LogInit()
	logger.Info("Логгер инициализирован")
	
	logger.Info("Проверка существования базы данных")
	db.DbExistCheck()
	defer func() {
		logger.Info("Отключение от базы данных")
		db.DB.Close()
	}()

	logger.Info("База данных успешно проверена и подключена")

	http.HandleFunc("/", handler.Handler)
	logger.Info("Маршруты HTTP зарегистрированы")

	logger.Info(fmt.Sprintf("Сервер запускается на %s", conf.Cfg.Host))
	err := http.ListenAndServe(conf.Cfg.Host, nil)
	if Error.GetErr(err) {
		logger.Fatal(fmt.Sprint(err))
	}
}
