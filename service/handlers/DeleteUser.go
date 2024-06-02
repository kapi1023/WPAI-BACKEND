package handlers

import (
	"log"
	"net/http"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func DeleteUserHandler(db *sqlx.DB) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		id := r.URL.Query().Get("id")

		var count int
		err := db.Get(&count, "SELECT COUNT(*) FROM reservations WHERE user_id=$1 AND end_date >= $2", id, time.Now())
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if count > 0 {
			http.Error(w, "Cannot delete a user with future reservations", http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM usersAirplane WHERE id=$1", id)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
