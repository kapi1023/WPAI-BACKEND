package models

import "time"

type Rental struct {
	ID         int       `json:"id"`
	UserID     int       `json:"userId"`
	AirplaneID int       `json:"airplaneId"`
	RentDate   time.Time `json:"rentDate"`
	Days       int       `json:"days"`
	Model      *string   `json:"model,omitempty"`
}
