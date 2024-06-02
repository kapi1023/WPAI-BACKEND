package models

type Airplane struct {
	ID           int     `json:"id" db:"id"`
	Model        string  `json:"model" db:"model"`
	Capacity     int     `json:"capacity" db:"capacity"`
	Availability bool    `json:"availability" db:"availability"`
	PricePerDay  float64 `json:"pricePerDay" db:"price_per_day"`
	TopSpeed     int     `json:"topSpeed" db:"top_speed"`
	FuelUsage    float64 `json:"fuelUsage" db:"fuel_usage"`
	Image        string  `json:"image" db:"image"`
}
