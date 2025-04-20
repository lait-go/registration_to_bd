package handler

import (
	"back/db"
	Error "back/internal/err"
	LogWork "back/internal/utils/log"
	"encoding/json"
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
	logger := LogWork.LogInit()
	path := strings.Trim(r.URL.Path, "/")

	if id, err := strconv.Atoi(path); err == nil {
		switch r.Method {
		case http.MethodPut:
			logger.Info(fmt.Sprintf("запрос на обновление пользователя %d получен", id))
			HandlerPUT(w, r, id)
		case http.MethodDelete:
			logger.Info(fmt.Sprintf("запрос на удаление пользователя %d получен", id))
			HandlerDELETE(w, r, id)
		case http.MethodGet:
			logger.Info(fmt.Sprintf("запрос на получение пользователя %d получен", id))
			HandlerGETByID(w, r, id)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
		return
	}

	switch r.Method {
	case http.MethodPost:
		logger.Info("запрос на добавление данных получен")
		HandlerPOST(w, r)
	case http.MethodGet:
		logger.Info("запрос на получение данных получен")
		HandlerGET(w, r, path)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

func HandlerPOST(w http.ResponseWriter, r *http.Request) {
	logger := LogWork.LogInit()
	w.Header().Set("Content-Type", "application/json")

	var date Post
	if err := json.NewDecoder(r.Body).Decode(&date); err != nil {
		Error.GetErr(err)
		http.Error(w, "Ошибка разбора тела запроса", http.StatusBadRequest)
		logger.Debug("не удалось декодировать тело запроса")
		return
	}

	logger.Debug("начато обогащение данных")

	date.Age = SafeGetAge(date.Name)
	if date.Age == 0 {logger.Debug("ошибка при получении возраста"); return}
	date.Gender = SafeGetGender(date.Name)
	if date.Gender == "" {logger.Debug("ошибка при получении пола"); return}
	date.Nationality = SafeGetNationality(date.Name)
	if date.Nationality == "" {logger.Debug("ошибка при получении национальности");return}

	logger.Debug("обогащение завершено")

	validate := validator.New()
	if err := validate.Struct(&date); err != nil {
		Error.GetErr(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Debug("валидация данных не пройдена")
		return
	}

	query := db.DbExecutor("../db/migrations/insert_to_tables.sql")
	_, err := db.DB.Exec(query, date.Name, date.Surname, date.Patronymic, date.Age, date.Gender, date.Nationality)
	if !Error.GetErr(err) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Пользователь успешно принят",
		})
		logger.Info("данные успешно записаны в бд")
	}
}

func HandlerGETByID(w http.ResponseWriter, r *http.Request, id int) {
	logger := LogWork.LogInit()
	query := db.DbExecutor("../db/migrations/receiving_by_id.sql")

	var date Post
	err := db.DB.Get(&date, query, id)
	if err != nil {
		logger.Debug("ошибка при получении данных по id")
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(date)
	logger.Info(fmt.Sprintf("данные пользователя %d успешно отправлены", id))
}

func HandlerGET(w http.ResponseWriter, r *http.Request, path string) {
	logger := LogWork.LogInit()

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

	logger.Info(fmt.Sprintf("тело запроса: %s", get))

	var date []Post
	err := db.DB.Select(&date, get)
	if err != nil {
		Error.GetErr(err)
		http.Error(w, "Ошибка при получении данных", http.StatusInternalServerError)
		logger.Debug("ошибка выполнения SELECT-запроса")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(date)
	logger.Info("данные успешно отправлены пользователю")
}