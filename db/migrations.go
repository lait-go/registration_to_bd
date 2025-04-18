package db

import (
	// conf "back/config"
	Error "back/internal/err"
	"back/internal/utils"
	"database/sql"
	"io"
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

var connStr = "postgres://myser:123@10.6.170.120:5432/mydb?sslmode=disable"


func initDbStruct() *sqlx.DB{
	db, err := sqlx.Open("postgres", connStr)
	Error.GetErr(err)

	err = db.Ping()
	Error.GetErr(err)

	return db
}

func DbExecutor(path string) sql.Result{
	db := initDbStruct()
	defer db.Close()

	file, err := utils.FileEx(path)
	if err != nil {
		Error.GetErr(err)
		log.Fatal()
	}
	
	defer file.Close()

	sqlDate, err := io.ReadAll(file)
	Error.GetErr(err)
	
	sqlRes := string(sqlDate)

	res, err := db.Exec(sqlRes)
	Error.GetErr(err)

	return res
}