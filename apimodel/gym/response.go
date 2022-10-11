package apimodel

import (
	"github.com/scalent-io/healthapi/apimodel"
	"github.com/scalent-io/healthapi/entity"
)

type CreateGymResponse struct {
	apimodel.Response
	Data string `json:"data,omitempty"`
}

type Gym struct {
	Name      string  `json:"name"`
	GymType   string  `json:"gym_type"`
	City      string  `json:"city"`
	State     string  `json:"state"`
	Country   string  `json:"country"`
	Address   string  `json:"address"`
	Lat       float64 `json:"lat"`
	Long      float64 `json:"long"`
	Amenities string  `json:"amenities"`
}

type GetAllGymResponse struct {
	apimodel.Response
	ImageBaseURL string `json:"imageBaseURL"`
	Page         int    `json:"page"`
	Limit        int    `json:"limit"`
	TotalRecords int    `json:"totalRecords"`
}

type GymDetailsResponse struct {
	ImageBaseURL string             `json:"imageBaseURL"`
	Gym          *entity.Gym        `json:"gym"`
	GymImages    []entity.GymImages `json:"gymImages"`
}

type GymSearchResponse struct {
	ImageBaseURL string `json:"imageBaseURL"`
	GymDistance  []entity.GymDistance
}
