package models

// CreateAdvertisementRequest параметры создания объявления; Price и Quantity не передаются, если nil
type CreateAdvertisementRequest struct {
	Title       string
	Description string
	Price       *float64
	Quantity    *int
	PhotoPaths  []string
}
