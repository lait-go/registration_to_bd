package db

import (
	conf "back/config"
	Error "back/internal/err"
	LogWork "back/internal/utils/log"

	"fmt"

	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func DbExistCheck(){
	logger := LogWork.LogInit()

	var err error
	DB, err = sqlx.Open("postgres", conf.Cfg.DateConf)
	if Error.GetErr(err) {
			logger.Fatal(fmt.Sprintf("Ошибка подключения к БД: %v", err))
	}

	if err = DB.Ping(); Error.GetErr(err) {
			logger.Fatal(fmt.Sprintf("БД недоступна: %v", err))
	}

	if DbExecutorNorParam("../db/migrations/check_tables.sql") == nil {
		DbExecutorNorParam("../db/migrations/create_tables.sql")
	}
}