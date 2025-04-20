package db

import (
	Error "back/internal/err"
	"back/internal/utils"
	LogWork "back/internal/utils/log"

	"database/sql"
	"io"

	_ "github.com/lib/pq" // PostgreSQL драйвер
)

func DbExecutor(path string) string {
	logger := LogWork.LogInit()
	logger.Info("Чтение SQL-файла: " + path)

	file, err := utils.FileEx(path)
	if Error.GetErr(err) {
		logger.Fatal("Не удалось открыть SQL-файл: " + path)
	}
	defer file.Close()

	sqlData, err := io.ReadAll(file)
	if Error.GetErr(err) {
		logger.Fatal("Не удалось прочитать SQL-файл: " + path)
	}

	logger.Info("SQL-файл успешно прочитан")
	return string(sqlData)
}

func DbExecutorNorParam(path string) sql.Result {
	logger := LogWork.LogInit()
	logger.Info("Чтение SQL-файла без параметров: " + path)

	file, err := utils.FileEx(path)
	if err != nil {
		Error.GetErr(err)
		logger.Fatal("Не удалось открыть SQL-файл: " + path)
	}
	defer file.Close()

	sqlData, err := io.ReadAll(file)
	if Error.GetErr(err) {
		logger.Fatal("Не удалось прочитать SQL-файл: " + path)
	}

	query := string(sqlData)
	logger.Info("Выполнение SQL-запроса из файла")

	res, err := DB.Exec(query)
	if Error.GetErr(err) {
		logger.Fatal("Ошибка при выполнении SQL-запроса")
	}

	logger.Info("SQL-запрос выполнен успешно")
	return res
}
