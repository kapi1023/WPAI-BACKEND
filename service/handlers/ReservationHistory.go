package handlers

import (
	"airplane/service/middleware"
	"airplane/service/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func ReservationHistoryHandler(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, ok := r.Context().Value(middleware.UserKey).(models.User)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query("SELECT r.id, a.model, r.start_date, r.end_date FROM reservations r JOIN airplanes a ON r.airplane_id = a.id WHERE r.user_id = $1", user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Println(err)
			return
		}
		defer rows.Close()

		var reservations []models.Reservation
		for rows.Next() {
			var reservation models.Reservation
			err := rows.Scan(&reservation.ID, &reservation.Model, &reservation.StartDate, &reservation.EndDate)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reservations = append(reservations, reservation)
		}
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(reservations)
	}
}
