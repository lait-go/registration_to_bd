package db

import (
	// conf "back/config"
	Error "back/internal/err"
	// "back/internal/utils"

	// "errors"
	"log"

	"github.com/jmoiron/sqlx"
)

var connStr = "postgres://myser:123@10.6.170.120:5432/mydb?sslmode=disable"
var DB *sqlx.DB

func DbExistCheck()  {
	var err error
	DB, err = sqlx.Open("postgres", connStr)
	if Error.GetErr(err) {
			log.Fatalf("Ошибка подключения к БД: %v", err)
	}

	if err = DB.Ping(); Error.GetErr(err) {
			log.Fatalf("БД недоступна: %v", err)
	}

	log.Println("база данных проверена")

	if DbExecutorNorParam("../db/migrations/check_tables.sql") == nil {
		DbExecutorNorParam("../db/migrations/create_tables.sql")
	}
	
}