package models

type Airplane struct {
	ID           int
	Model        string
	Capacity     int
	Availability bool
	PricePerDay  float64
	TopSpeed     int
	FuelUsage    float64
	Image        string
}
