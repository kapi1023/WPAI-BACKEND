package handlers

import (
	"airplane/service/models"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func AirplanesHandler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

		rows, err := db.Query("SELECT id, model, capacity, availability, price_per_day FROM airplanes")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var airplanes []models.Airplane
		for rows.Next() {
			var airplane models.Airplane
			if err := rows.Scan(&airplane.ID, &airplane.Model, &airplane.Capacity, &airplane.Availability, &airplane.PricePerDay); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			airplanes = append(airplanes, airplane)
		}

		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(airplanes)
	}
}
