package model

import "math"

type Location struct {
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (l *Location) Valid() bool {
	return ((l.Latitude >= -90) && (l.Latitude <= 90)) && ((l.Longitude >= -180) && (l.Longitude <= 180))
}

func (l *Location) Distance(other Location) float64 {
	// Haversine formula
	const r float64 = 6371
	lat2 := other.Latitude * math.Pi / 180
	lat1 := l.Latitude * math.Pi / 180
	long2 := other.Longitude * math.Pi / 180
	long1 := l.Longitude * math.Pi / 180
	dLat := (lat2 - lat1) / 2
	dLong := (long2 - long1) / 2
	return 2 * r * math.Asin(math.Sqrt(math.Pow(math.Sin(dLat), 2) + (math.Cos(lat1) * math.Cos(lat2) * math.Pow(math.Sin(dLong), 2))))
}
