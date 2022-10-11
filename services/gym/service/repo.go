package service

import (
	"context"

	"github.com/scalent-io/healthapi/entity"
	"github.com/scalent-io/healthapi/pkg/errors"
)

//go:generate mockgen -source=repo.go -destination=repomock/repo_mock.go -package=repomock

type GymRepo interface {
	Create(ctx context.Context, gym *entity.Gym) errors.Response
	Update(ctx context.Context, gym *entity.Gym) errors.Response
	GetAll(ctx context.Context, limit, offset int) ([]entity.Gym, errors.Response)
	Count(ctx context.Context) (int, errors.Response)
	GetById(ctx context.Context, gymId int) (*entity.Gym, errors.Response)
	Search(ctx context.Context, lat, long float64) ([]entity.GymDistance, errors.Response)
	GetGymImages(ctx context.Context, limit, offset int) ([]entity.GymDetails, errors.Response)
	CountGetGymImages(ctx context.Context) (int, errors.Response)
}

type GymImagesRepo interface {
	Create(ctx context.Context, gym entity.GymImages) errors.Response
	GetImagesId(ctx context.Context, gymId int) ([]entity.GymImages, errors.Response)
}
