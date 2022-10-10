package entity

import "time"

type Gym struct {
	ID         int
	Name       string
	GymType    string
	City       string
	State      string
	Country    string
	Address    string
	Lat        float64
	Long       float64
	Amenities  string
	CreatedOn  time.Time
	ModifiedOn time.Time
}

type GymDistance struct {
	ID         int
	Name       string
	GymType    string
	City       string
	State      string
	Country    string
	Address    string
	Lat        float64
	Long       float64
	Amenities  string
	Distance   float64
	CreatedOn  time.Time
	ModifiedOn time.Time
}

type GymImages struct {
	ID         int
	GymID      int
	ImageType  string
	Label      string
	CreatedOn  time.Time
	ModifiedOn time.Time
}

type GymDetails struct {
	ID         int
	Name       string
	GymType    string
	City       string
	State      string
	Country    string
	Address    string
	Lat        float64
	Long       float64
	Amenities  string
	CreatedOn  time.Time
	ModifiedOn time.Time
	ImageType  string
	Label      string
}

type GymCount struct {
	TotalRecords int `pg:"total_records"`
}
