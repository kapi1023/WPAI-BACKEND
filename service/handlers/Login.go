package handlers

import (
	"airplane/service/models"
	"airplane/service/utils"
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func LoginHandler(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var creds models.Credentials
		err := json.NewDecoder(r.Body).Decode(&creds)
		if err != nil {
			http.Error(w, "Invalid request payload", http.StatusBadRequest)
			return
		}

		var user models.User
		err = db.QueryRow("SELECT id, username, password, is_admin FROM usersAirplane WHERE username = $1", creds.Username).Scan(&user.ID, &user.Username, &user.Password, &user.IsAdmin)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		if !utils.CheckPasswordHash(creds.Password, user.Password) {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		token, err := utils.GenerateToken(user.ID, user.IsAdmin)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"token":    token,
			"username": user.Username,
			"isAdmin":  user.IsAdmin,
		})
	}
}
