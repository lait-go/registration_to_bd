package handler

import (
	"back/db"
	Error "back/internal/err"
	LogWork "back/internal/utils/log"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

func HandlerPUT(w http.ResponseWriter, r *http.Request, id int) {
	logger := LogWork.LogInit()
	w.Header().Set("Content-Type", "application/json")

	var updated Post
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		Error.GetErr(err)
		http.Error(w, "Ошибка разбора тела запроса", http.StatusBadRequest)
		logger.Debug("не удалось декодировать тело запроса в PUT")
		return
	}

	updated.ID = id

	logger.Debug("начато обогащение данных для обновления")
	updated.Age = SafeGetAge(updated.Name)
	if updated.Age == 0 {
		logger.Debug("ошибка при получении возраста")
		return
	}
	updated.Gender = SafeGetGender(updated.Name)
	if updated.Gender == "" {
		logger.Debug("ошибка при получении пола")
		return
	}
	updated.Nationality = SafeGetNationality(updated.Name)
	if updated.Nationality == "" {
		logger.Debug("ошибка при получении национальности")
		return
	}
	logger.Debug("обогащение завершено")

	validate := validator.New()
	if err := validate.Struct(&updated); err != nil {
		Error.GetErr(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Debug("валидация обновлённых данных не пройдена")
		return
	}

	query := db.DbExecutor("../db/migrations/update_by_id.sql")
	_, err := db.DB.Exec(query, updated.Name, updated.Surname, updated.Patronymic, updated.Age, updated.Gender, updated.Nationality, updated.ID)
	if !Error.GetErr(err) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": fmt.Sprintf("Пользователь с id=%d успешно обновлён", updated.ID),
		})
		logger.Info(fmt.Sprintf("данные пользователя %d успешно обновлены", updated.ID))
	}
}

func HandlerDELETE(w http.ResponseWriter, r *http.Request, id int) {
	logger := LogWork.LogInit()

	query := db.DbExecutor("../db/migrations/delete_by_id.sql")

	res, err := db.DB.Exec(query, id)
	if err != nil {
		Error.GetErr(err)
		http.Error(w, "Ошибка при удалении данных", http.StatusInternalServerError)
		logger.Debug("ошибка при удалении данных из БД")
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		Error.GetErr(err)
		http.Error(w, "Ошибка при проверке удаления", http.StatusInternalServerError)
		logger.Debug("ошибка при вызове RowsAffected")
		return
	}

	if rows == 0 {
		http.Error(w, "Пользователь с таким id не найден", http.StatusNotFound)
		logger.Debug(fmt.Sprintf("пользователь с id=%d не найден", id))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": fmt.Sprintf("Пользователь c id=%d успешно удалён", id),
	})
	logger.Info(fmt.Sprintf("пользователь с id=%d успешно удалён", id))
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
