package converter

import (
	gymApimodel "github.com/scalent-io/healthapi/apimodel/gym"
	"github.com/scalent-io/healthapi/entity"
)

func GymAPIToGymEntity(m gymApimodel.CreateGymReq) entity.Gym {
	var e entity.Gym

	e.Name = m.Name
	e.GymType = m.GymType
	e.City = m.City
	e.State = m.State
	e.Country = m.Country
	e.Address = m.Address
	e.Amenities = m.Amenities
	e.Lat = m.Lat
	e.Long = m.Long
	return e
}

func GymAPIToGymImagesEntity(images []gymApimodel.CreateGymImageReq) []entity.GymImages {
	var e []entity.GymImages

	for _, img := range images {
		entity := entity.GymImages{
			GymID:     img.GymID,
			ImageType: img.ImageType,
			Label:     img.Label,
			Type:      img.Type,
		}
		e = append(e, entity)
	}
	return e
}

func GymAPIToLogoEntity(m gymApimodel.CreateLogoReq) entity.Gym {
	var e entity.Gym
	e.ID = m.GymID
	e.LogoName = m.LogoName
	return e
}
