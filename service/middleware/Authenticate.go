package middleware

import (
	"airplane/service/models"
	"airplane/service/utils"
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

type key string

var jwtKey = []byte("5zLZHHndhY95lbJPgOweh4dKRbOeH8c4")

const UserKey key = "user"

func Auth(db *sqlx.DB) func(httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			claims := &utils.Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(jwtKey), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			var user models.User
			err = db.QueryRow("SELECT id, username, is_admin FROM usersAirplane WHERE id = $1", claims.UserID).Scan(&user.ID, &user.Username, &user.IsAdmin)
			if err != nil {
				http.Error(w, "User not found", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, user)
			next(w, r.WithContext(ctx), ps)
		}
	}
}

func Admin() func(httprouter.Handle) httprouter.Handle {
	return func(next httprouter.Handle) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
				return
			}

			tokenString := parts[1]
			claims := &utils.Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if !claims.IsAdmin {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next(w, r, ps)
		}
	}
}
