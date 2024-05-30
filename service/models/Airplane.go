package models

type Airplane struct {
	ID           int     `json:"id"`
	Model        string  `json:"model"`
	Capacity     int     `json:"capacity"`
	Availability bool    `json:"availability"`
	PricePerDay  float64 `json:"pricePerDay"`
	TopSpeed     int     `json:"topSpeed"`
	FuelUsage    float64 `json:"fuelUsage"`
	Image        string  `json:"image"`
}
