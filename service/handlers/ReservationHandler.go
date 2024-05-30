package handlers

import (
	"airplane/service/models"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func ReservationHistoryHandler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		user, ok := r.Context().Value("user").(models.User)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		rows, err := db.Query("SELECT r.id, a.model, r.rent_date, r.days FROM rentals r JOIN airplanes a ON r.airplane_id = a.id WHERE r.user_id = $1", user.ID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var reservations []models.Rental
		for rows.Next() {
			var reservation models.Rental
			err := rows.Scan(&reservation.ID, &reservation.Model, &reservation.RentDate, &reservation.Days)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			reservations = append(reservations, reservation)
		}

		json.NewEncoder(w).Encode(reservations)
	}
}
