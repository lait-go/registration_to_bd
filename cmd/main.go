package main

import (
	conf "back/config"
	"back/db"
	handler "back/internal/api"
	"back/internal/utils/log"
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

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	path := r.URL.Path
	// 	if path == "/" {
	// 		http.ServeFile(w, r, "./../static/index.html")
	// 		return
	// 	}
	// 	http.ServeFile(w, r, "./../static"+path)
	// })
	
	log.Println("Сервер запущен на http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}
