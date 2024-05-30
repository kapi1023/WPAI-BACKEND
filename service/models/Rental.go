package models

import "time"

type Rental struct {
	ID         int
	UserID     int
	AirplaneID int
	RentDate   time.Time
	Days       int
}
