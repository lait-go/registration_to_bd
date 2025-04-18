package handler

import (
	Error "back/internal/err"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

type Post struct {
	Name 			 string   `json:"name" validate:"required"`
	Surname    string		`json:"surname" validate:"required"`
	Patronymic string   `json:"patronymic"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost{
		http.Error(w, "Не верный метод", http.StatusMethodNotAllowed)
		var err = errors.New("не верный метод")
		Error.GetErr(err)
		return
	}

	var date Post

	err := json.NewDecoder(r.Body).Decode(&date)
	Error.GetErr(err)

  validate := validator.New()
	err = validate.Struct(&date)
	Error.GetErr(err)
	
	fmt.Println(date)
}

func HandlerPOST(w http.ResponseWriter, r *http.Request){

}