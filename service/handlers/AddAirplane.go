package handlers

import (
	"airplane/service/models"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func AddAirplaneHandler(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var airplane models.Airplane
		err := json.NewDecoder(r.Body).Decode(&airplane)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = db.Exec("INSERT INTO airplanes (model, capacity, availability, price_per_day, top_speed, fuel_usage, image) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			airplane.Model, airplane.Capacity, true, airplane.PricePerDay, airplane.TopSpeed, airplane.FuelUsage, airplane.Image)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
