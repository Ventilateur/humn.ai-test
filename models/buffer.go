package models

type Coordinate struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

type Input struct {
	Coordinate
}

type Output struct {
	Coordinate
	Postcode string `json:"postcode"`
}
