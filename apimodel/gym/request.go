package apimodel

type CreateGymReq struct {
	Name      string  `json:"name" validate:"required"`
	GymType   string  `json:"gym_type" validate:"required"`
	City      string  `json:"city" validate:"required"`
	State     string  `json:"state" validate:"required"`
	Country   string  `json:"country" validate:"required"`
	Address   string  `json:"address" validate:"required"`
	Lat       float64 `json:"lat" validate:"required" `
	Long      float64 `json:"long" validate:"required"`
	Amenities string  `json:"amenities"`
}

type CreateGymImageReq struct {
	GymID     int    `json:"gym_id"`
	ImageType string `json:"image_type"`
	Label     string `json:"label"`
	Type      string `json:"type"`
}

type CreateLogoReq struct {
	GymID    int    `json:"gym_id"`
	LogoName string `json:"logo_name"`
}
