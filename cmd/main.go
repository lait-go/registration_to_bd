package main

import (
	conf "back/config"
	"back/db"
	handler "back/internal/api"
	LogWork "back/internal/utils/log"
	"fmt"
	"log"

	"net/http"
)

func main() {
	conf.Cfg = conf.Config_work()
	logger := LogWork.LogInit()
	logger.Info("конфигурация считана")
	logger.Info("логер запущен")
	db.DbExistCheck()
	logger.Info("база данных проверена")

	
	http.HandleFunc("/", handler.Handler)
	
	logger.Info(fmt.Sprintf("Сервер запущен на %s", conf.Cfg.Address))
	err := http.ListenAndServe(conf.Cfg.Address, nil)
	if err != nil {
		log.Fatal(err)
	}
}
