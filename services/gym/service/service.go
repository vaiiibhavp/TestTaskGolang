package service

import (
	"context"
	"net/url"

	"github.com/scalent-io/healthapi/apimodel"
	gymapimodel "github.com/scalent-io/healthapi/apimodel/gym"
	"github.com/scalent-io/healthapi/entity"
	"github.com/scalent-io/healthapi/pkg/errors"
	"github.com/scalent-io/healthapi/pkg/server"
)

//go:generate mockgen -source=service.go -destination=servicemock/service_mock.go -package=servicemock

type GymService interface {
	Create(ctx context.Context, gym *entity.Gym) (*apimodel.Response, errors.Response)
	GetAll(ctx context.Context, params url.Values, imageConfig *server.ImageConfig) (*gymapimodel.GetAllGymResponse, errors.Response)
	GetById(ctx context.Context, gymID int, imageConfig *server.ImageConfig) (*apimodel.Response, errors.Response)
	Search(ctx context.Context, params url.Values, imageConfig *server.ImageConfig) (*apimodel.Response, errors.Response)
	Upload(ctx context.Context, images []entity.GymImages) (*apimodel.Response, errors.Response)
	UploadLogo(ctx context.Context, req gymapimodel.CreateLogoReq) (*apimodel.Response, errors.Response)
}
