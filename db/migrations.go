package db

import (
	Error "back/internal/err"
	"back/internal/utils"
	"database/sql"
	"io"
	"log"

	"github.com/jmoiron/sqlx"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

func InitDbStruct() *sqlx.DB{
	db, err := sqlx.Open("postgres", connStr)
	Error.GetErr(err)

	err = db.Ping()
	Error.GetErr(err)

	return db
}

func DbExecutor(path string) string {
	file, err := utils.FileEx(path)
	Error.GetErr(err)

	sqlDate, err := io.ReadAll(file)
	Error.GetErr(err)

	return string(sqlDate)
}

func DbExecutorNorParam(path string) sql.Result{
	db := InitDbStruct()

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