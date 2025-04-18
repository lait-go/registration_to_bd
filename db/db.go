package db

import (
	conf "back/config"
	Error "back/internal/err"
	"back/internal/utils"
	// "errors"
	"log"
)

func DbExistCheck()  {
	file, err := utils.FileEx(conf.Cfg.Path.StoragePath)
	if err != nil {
		Error.GetErr(err)
		log.Fatal()
	}
	
	defer file.Close()

	if DbExecutor("../db/migrations/check_tables.sql") == nil {
		DbExecutor("../db/migrations/create_tables.sql")
		DbExecutor("../db/migrations/add_constraints_to_tables.sql")
	}
}