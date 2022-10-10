package service

import (
	"context"
	"net/url"

	"github.com/scalent-io/healthapi/apimodel"
	gymapimodel "github.com/scalent-io/healthapi/apimodel/gym"
	"github.com/scalent-io/healthapi/entity"
	"github.com/scalent-io/healthapi/pkg/errors"
)

//go:generate mockgen -source=service.go -destination=servicemock/service_mock.go -package=servicemock

type GymService interface {
	Create(ctx context.Context, gym *entity.Gym) (*apimodel.Response, errors.Response)
	GetAll(ctx context.Context, params url.Values) (*gymapimodel.GetAllGymResponse, errors.Response)
	GetById(ctx context.Context, gymID int) (*apimodel.Response, errors.Response)
	Search(ctx context.Context, params url.Values) (*apimodel.Response, errors.Response)
	Upload(ctx context.Context, images []entity.GymImages) (*apimodel.Response, errors.Response)
}
