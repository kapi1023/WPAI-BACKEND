package service

import (
	"airplane/service/handlers"
	"airplane/service/middleware"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/julienschmidt/httprouter"
)

func Start(db *sqlx.DB) {
	router := httprouter.New()

	router.GET("/airplanes", handlers.AirplanesHandler(db))
	router.GET("/airplane", handlers.AirplaneDetailsHandler(db))

	router.GET("/reservation_history", middleware.Auth(db)(handlers.ReservationHistoryHandler(db)))
	router.POST("/rent", middleware.Auth(db)(handlers.RentHandler(db)))

	router.POST("/register", handlers.RegisterHandler(db))
	router.POST("/login", handlers.LoginHandler(db))

	router.POST("/addAirplane", middleware.Admin()(handlers.AddAirplaneHandler(db)))
	router.PUT("/editAirplane", middleware.Admin()(handlers.EditAirplaneHandler(db)))
	router.DELETE("/deleteAirplane", middleware.Admin()(handlers.DeleteAirplaneHandler(db)))
	router.DELETE("/deleteRent", middleware.Admin()(handlers.DeleteRentHandler(db)))
	router.GET("/reservations", middleware.Admin()(handlers.ListRentsHandler(db)))

	log.Println("Server started at :80")
	log.Fatal(http.ListenAndServe(":80", middleware.Cors(router)))
}
