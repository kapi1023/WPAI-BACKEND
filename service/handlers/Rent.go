package handlers

import (
	"airplane/service/models"
	"airplane/service/utils"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

func RentHandler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var rental models.Rental
		err := json.NewDecoder(r.Body).Decode(&rental)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tokenCookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID, err := utils.GetUserIDFromToken(tokenCookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rental.UserID = userID
		rental.RentDate = time.Now()

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var availability bool
		err = tx.QueryRow("SELECT availability FROM airplanes WHERE id=$1 FOR UPDATE", rental.AirplaneID).Scan(&availability)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				http.Error(w, "Airplane not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if !availability {
			tx.Rollback()
			http.Error(w, "Airplane not available", http.StatusBadRequest)
			return
		}

		_, err = tx.Exec("UPDATE airplanes SET availability=false WHERE id=$1", rental.AirplaneID)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err = tx.Exec("INSERT INTO rentals (user_id, airplane_id, rent_date, days) VALUES ($1, $2, $3, $4)", rental.UserID, rental.AirplaneID, rental.RentDate, rental.Days)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tx.Commit()
		w.WriteHeader(http.StatusCreated)
	}
}
