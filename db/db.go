package db

import (
	conf "back/config"
	Error "back/internal/err"
	LogWork "back/internal/utils/log"

	"fmt"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func DbExistCheck() {
	logger := LogWork.LogInit()

	logger.Info("Инициализация подключения к базе данных")
	var err error
	DB, err = sqlx.Open("postgres", conf.Cfg.DateConf)
	if Error.GetErr(err) {
		logger.Fatal(fmt.Sprintf("Ошибка при открытии соединения с БД: %v", err))
	}

	logger.Info("Проверка доступности базы данных (ping)")
	if err = DB.Ping(); Error.GetErr(err) {
		logger.Fatal(fmt.Sprintf("БД недоступна: %v", err))
	}
	logger.Info("Соединение с базой данных успешно установлено")

	logger.Info("Проверка существования таблиц")
	if DbExecutorNorParam("../db/migrations/check_tables.sql") == nil {
		logger.Info("Таблицы не найдены, выполняется создание")
		DbExecutorNorParam("../db/migrations/create_tables.sql")
		logger.Info("Таблицы успешно созданы")
	} else {
		logger.Info("Таблицы уже существуют")
	}
}
