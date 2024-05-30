package service

import (
	"airplane/service/handlers"
	"airplane/service/middleware"
	"database/sql"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Start(db *sql.DB) {
	router := httprouter.New()

	router.POST("/register", handlers.RegisterHandler(db))
	router.POST("/login", handlers.LoginHandler(db))
	router.GET("/airplanes", handlers.AirplanesHandler(db))
	router.GET("/airplanes/:id", handlers.AirplaneDetailsHandler(db))
	router.POST("/rent", middleware.Auth(handlers.RentHandler(db)))
	router.POST("/add_airplane", middleware.Admin(handlers.AddAirplaneHandler(db)))
	router.POST("/edit_airplane", middleware.Admin(handlers.EditAirplaneHandler(db)))
	router.POST("/delete_airplane", middleware.Admin(handlers.DeleteAirplaneHandler(db)))

	log.Fatal(http.ListenAndServe(":8080", router))
}
