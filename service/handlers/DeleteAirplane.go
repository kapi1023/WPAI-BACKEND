package handlers

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func DeleteAirplaneHandler(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		id := r.URL.Query().Get("id")

		tx, err := db.Begin()
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var count int
		err = tx.QueryRow("SELECT COUNT(*) FROM reservations WHERE airplane_id = $1 AND end_date >= $2", id, time.Now()).Scan(&count)
		if err != nil {
			tx.Rollback()
			if err == sql.ErrNoRows {
				log.Println(err)
				http.Error(w, "Airplane not found", http.StatusNotFound)
			} else {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		if count > 0 {
			tx.Rollback()
			log.Println("Cannot delete an airplane with future reservations")
			http.Error(w, "Cannot delete an airplane with future reservations", http.StatusBadRequest)
			return
		}

		_, err = tx.Exec("DELETE FROM airplanes WHERE id=$1", id)
		if err != nil {
			tx.Rollback()
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tx.Commit()

		w.WriteHeader(http.StatusOK)
	}
}
