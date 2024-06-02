package handlers

import (
	"airplane/service/middleware"
	"airplane/service/models"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

type RentRequest struct {
	AirplaneID int    `json:"airplane_id"`
	StartDate  string `json:"start_date"`
	EndDate    string `json:"end_date"`
}

func RentHandler(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var req RentRequest
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, ok := r.Context().Value(middleware.UserKey).(models.User)
		if !ok {
			log.Println(err)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		startDate, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			http.Error(w, "Invalid start date format", http.StatusBadRequest)
			return
		}

		endDate, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			http.Error(w, "Invalid end date format", http.StatusBadRequest)
			return
		}

		if startDate.Before(time.Now()) {
			log.Println(err)
			http.Error(w, "Start date must be in the future", http.StatusBadRequest)
			return
		}

		if endDate.Before(startDate) {
			log.Println(err)
			http.Error(w, "Invalid date range", http.StatusBadRequest)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var count int
		err = tx.QueryRow("SELECT COUNT(*) FROM reservations WHERE airplane_id = $1 AND (start_date <= $2 AND end_date >= $3)", req.AirplaneID, req.EndDate, req.StartDate).Scan(&count)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count > 0 {
			tx.Rollback()
			log.Println(err)
			http.Error(w, "Airplane not available for the selected period", http.StatusBadRequest)
			return
		}

		_, err = tx.Exec("INSERT INTO reservations (user_id, airplane_id, start_date, end_date) VALUES ($1, $2, $3, $4)",
			user.ID, req.AirplaneID, req.StartDate, req.EndDate)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tx.Commit()

		w.WriteHeader(http.StatusCreated)
	}
}
