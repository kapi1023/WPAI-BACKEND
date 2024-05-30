package handlers

import (
	"airplane/service/models"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func DeleteAirplaneHandler(db *sql.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		var airplane models.Airplane
		err := json.NewDecoder(r.Body).Decode(&airplane)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		tx, err := db.Begin()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var availability bool
		err = tx.QueryRow("SELECT availability FROM airplanes WHERE id=$1 FOR UPDATE", airplane.ID).Scan(&availability)
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
			http.Error(w, "Cannot delete an airplane that is rented", http.StatusBadRequest)
			return
		}

		_, err = tx.Exec("DELETE FROM airplanes WHERE id=$1", airplane.ID)
		if err != nil {
			tx.Rollback()
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tx.Commit()
		w.WriteHeader(http.StatusOK)
	}
}
