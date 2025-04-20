package handler

import (
	"back/db"
	Error "back/internal/err"
	LogWork "back/internal/utils/log"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func HandleId(w http.ResponseWriter, r *http.Request, id int) {
	logger := LogWork.LogInit()

	dateStr := db.DbExecutor("../db/migrations/delete_by_id.sql")

	res,err := db.DB.Exec(dateStr, id)
	if Error.GetErr(err){
		http.Error(w, "Ошибка при удалении данных", http.StatusInternalServerError)
		logger.Debug("ошибка при удалении данных")
		return
	}

	if ready, err := res.RowsAffected(); ready == 0 || Error.GetErr(err){
		http.Error(w, "пользователь с таким id не найден", http.StatusInternalServerError)
		logger.Debug("неверный id: HandlerID")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Пользователь c id=%d успешно удален", id),
	})
	logger.Info(fmt.Sprintf("данные пользователя %d успешно удалены", id))
}

func plag(path, query, key string, isString bool) string {
	if strings.Contains(path, key+"=") {
		parts := strings.Split(path, key+"=")
		if len(parts) >= 2 {
			value := strings.SplitN(parts[1], "&", 2)[0]

			switch key {
			case "age_min":
				query += fmt.Sprintf(" AND AGE >= %s", value)
			case "age_max":
				query += fmt.Sprintf(" AND AGE < %s", value)
			case "limit", "offset":
				query += fmt.Sprintf(" %s %s", strings.ToUpper(key), value)
			default:
				if isString {
					query += fmt.Sprintf(" AND %s = '%s'", strings.ToUpper(key), value)
				} else {
					query += fmt.Sprintf(" AND %s = %s", strings.ToUpper(key), value)
				}
			}
		}
	}
	return query
}
