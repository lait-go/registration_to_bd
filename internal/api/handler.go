package handler

import (
	"back/db"
	Error "back/internal/err"
	LogWork "back/internal/utils/log"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Post struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name" validate:"required"`
	Surname     string `json:"surname" db:"surname" validate:"required"`
	Patronymic  string `json:"patronymic" db:"patronymic"`
	Age         int    `json:"age" db:"age" validate:"gte=1,lte=100"`
	Gender      string `json:"gender" db:"gender" validate:"oneof=male female"`
	Nationality string `json:"nationality" db:"nationality" validate:"required"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	path := strings.Trim(r.URL.Path, "/")

	if id, err := strconv.Atoi(path); err == nil {
			HandleId(w, r, id)
	} else if path == "" && r.Method == http.MethodPost{
		HandlerPOST(w, r)
	}else{
		HandlerGET(w, r, path)
	}
}

func HandlerPOST(w http.ResponseWriter, r *http.Request){
	logger := LogWork.LogInit()
	logger.Info("Получен запрос(POST)")

	w.Header().Set("Content-Type", "application/json")

	var date Post

	err := json.NewDecoder(r.Body).Decode(&date)
	Error.GetErr(err)

	logger.Debug("запрос Age")
	date.Age = GetAge(date.Name)
	logger.Debug("запрос Gender")
	date.Gender = GetGender(date.Name)
	logger.Debug("запрос Nationality")
	date.Nationality = GetNationality(date.Name)

	validate := validator.New()
	if err = validate.Struct(&date); err != nil{
		Error.GetErr(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	logger.Debug("данные успешно приняты")

	dateStr := db.DbExecutor("../db/migrations/insert_to_tables.sql")
	
	_, err = db.DB.Exec(dateStr, date.Name, date.Surname, date.Patronymic, date.Age, date.Gender, date.Nationality)
	if !Error.GetErr(err){
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]string{
			"message": "Пользователь успешно принят",
		})
		
		logger.Info("данные успешно записаны в бд, пользователь получил ответ")
	}
}

func HandleId(w http.ResponseWriter, r *http.Request, id int) {
	var date Post

	dateStr := db.DbExecutor("../db/migrations/receiving_by_id.sql")

	err := db.DB.Get(&date, dateStr, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}
		Error.GetErr(err)
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(date)
}

func HandlerGET(w http.ResponseWriter, r *http.Request, path string) {
	get := "SELECT * FROM person WHERE 1=1"

	get = plag(path, get, "id", false)
	get = plag(path, get, "name", true)
	get = plag(path, get, "surname", true)
	get = plag(path, get, "patronymic", true)
	get = plag(path, get, "age", false)
	get = plag(path, get, "gender", true)
	get = plag(path, get, "nationality", true)
	get = plag(path, get, "age_min", false)
	get = plag(path, get, "age_max", false)
	get = plag(path, get, "limit", false)
	get = plag(path, get, "offset", false)

	var date []Post
	err := db.DB.Select(&date, get)
	if err != nil {
		Error.GetErr(err)
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(date)
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

