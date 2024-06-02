package handlers

import (
	"airplane/service/models"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func AirplanesHandler(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		startDateStr := r.URL.Query().Get("start_date")
		endDateStr := r.URL.Query().Get("end_date")

		var startDate, endDate time.Time
		var err error

		if startDateStr != "" && endDateStr != "" {
			startDate, err = time.Parse("2006-01-02", startDateStr)
			if err != nil {
				http.Error(w, "Invalid start date format", http.StatusBadRequest)
				return
			}
			endDate, err = time.Parse("2006-01-02", endDateStr)
			if err != nil {
				http.Error(w, "Invalid end date format", http.StatusBadRequest)
				return
			}
		}

		query := `
			SELECT a.id, a.model, a.capacity, a.availability, a.price_per_day, a.top_speed, a.fuel_usage, a.image
			FROM airplanes a
		`

		args := []interface{}{}
		if startDateStr != "" && endDateStr != "" {
			if endDate.Before(startDate) {
				http.Error(w, "Invalid date range", http.StatusBadRequest)
				return
			}

			query += `
				WHERE a.id NOT IN (
					SELECT r.airplane_id
					FROM reservations r
					WHERE $1 BETWEEN r.start_date AND r.end_date
					OR $2 BETWEEN r.start_date AND r.end_date
				)
			`
			args = append(args, startDate, endDate)
		}

		rows, err := db.Queryx(query, args...)
		if err != nil {
			http.Error(w, "Failed to query airplanes", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var airplanes []models.Airplane
		for rows.Next() {
			var airplane models.Airplane
			if err := rows.StructScan(&airplane); err != nil {
				http.Error(w, "Failed to scan airplane", http.StatusInternalServerError)
				return
			}
			airplanes = append(airplanes, airplane)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(airplanes)
	}
}
