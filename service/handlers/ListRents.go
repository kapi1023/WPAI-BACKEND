package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

type ListReservationResponse struct {
	ID        int    `json:"id"`
	Model     string `json:"model"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

func ListRentsHandler(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		rows, err := db.Query(`Select reservations.id, airplanes.model, reservations.start_date, reservations.end_date
		FROM reservations
		join airplanes on airplanes.Id = reservations.airplane_id
		`)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var rents []ListReservationResponse
		for rows.Next() {
			var rent ListReservationResponse
			if err := rows.Scan(&rent.ID, &rent.Model, &rent.StartDate, &rent.EndDate); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			rents = append(rents, rent)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(rents)
	}
}
