package model

import "time"

type Gym struct {
	ID         int       `pg:",pk"`
	Name       string    `pg:"name"`
	GymType    string    `pg:"gym_type"`
	City       string    `pg:"city"`
	State      string    `pg:"state"`
	Country    string    `pg:"country"`
	Address    string    `pg:"address"`
	Lat        float64   `pg:"lat"`
	Long       float64   `pg:"long"`
	Amenities  string    `pg:"amenities"`
	LogoName   string    `pg:"logo_name"`
	CreatedOn  time.Time `pg:"created_on"`
	ModifiedOn time.Time `pg:"modified_on"`
}
