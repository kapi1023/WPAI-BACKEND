package middleware

import (
	"airplane/service/utils"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Auth(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		_, err = utils.ParseToken(tokenCookie.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r, ps)
	}
}

func Admin(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		tokenCookie, err := r.Cookie("token")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		claims, err := utils.ParseToken(tokenCookie.Value)
		if err != nil || !claims.IsAdmin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next(w, r, ps)
	}
}
