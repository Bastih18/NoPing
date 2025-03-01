package globals

import "time"

type Stats struct {
	Attempted int
	Connected int
	Failed    int
	TotalTime time.Duration
	Min       time.Duration
	Max       time.Duration
}

type GeoInfo struct {
	City    string `json:"city"`
	Region  string `json:"region"`
	Country string `json:"country"`
}
