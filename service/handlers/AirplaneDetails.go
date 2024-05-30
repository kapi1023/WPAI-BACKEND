package handlers

import (
	"airplane/service/models"
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

func AirplaneDetailsHandler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id, err := strconv.Atoi(ps.ByName("id"))
		if err != nil {
			http.Error(w, "Invalid airplane ID", http.StatusBadRequest)
			return
		}

		var airplane models.Airplane
		err = db.QueryRow("SELECT id, model, capacity, availability, price_per_day, top_speed, fuel_usage, image FROM airplanes WHERE id=$1", id).Scan(
			&airplane.ID, &airplane.Model, &airplane.Capacity, &airplane.Availability, &airplane.PricePerDay, &airplane.TopSpeed, &airplane.FuelUsage, &airplane.Image,
		)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Airplane not found", http.StatusNotFound)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(airplane)
	}
}
