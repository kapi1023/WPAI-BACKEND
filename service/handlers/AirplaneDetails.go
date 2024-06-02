package handlers

import (
	"airplane/service/models"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func AirplaneDetailsHandler(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		var airplane models.Airplane
		var err error
		id := r.URL.Query().Get("id")

		err = db.QueryRow("SELECT id, model, capacity, availability, price_per_day, top_speed, fuel_usage, image FROM airplanes WHERE id=$1", id).Scan(
			&airplane.ID, &airplane.Model, &airplane.Capacity, &airplane.Availability, &airplane.PricePerDay, &airplane.TopSpeed, &airplane.FuelUsage, &airplane.Image,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Println(err)
				http.Error(w, "Airplane not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(airplane)
	}
}
