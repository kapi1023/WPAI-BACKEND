package handlers

import (
	"airplane/service/models"
	"airplane/service/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func RegisterHandler(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !utils.ValidatePassword(user.Password) {
			log.Println("Password must be at least 8 characters long and include an uppercase letter, a lowercase letter, a number, and a special character.")
			http.Error(w, "Password must be at least 8 characters long and include an uppercase letter, a lowercase letter, a number, and a special character.", http.StatusBadRequest)
			return
		}

		hashedPassword, err := utils.HashPassword(user.Password)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		user.Password = hashedPassword

		_, err = db.Exec("INSERT INTO usersAirplane (username, password, is_admin) VALUES ($1, $2, $3)", user.Username, user.Password, user.IsAdmin)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
