package models

import "time"

type Reservation struct {
	ID         int       `json:"id"`
	UserID     int       `json:"userId"`
	AirplaneID int       `json:"airplaneId"`
	RentDate   time.Time `json:"rentDate"`
	StartDate  time.Time `json:"startDate"`
	EndDate    time.Time `json:"endDate"`
	Model      *string   `json:"model,omitempty"`
}
