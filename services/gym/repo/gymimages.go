package repo

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/go-pg/pg/v10"

	"github.com/scalent-io/healthapi/entity"
	c "github.com/scalent-io/healthapi/pkg/context"
	"github.com/scalent-io/healthapi/pkg/errors"
)

type GymImageRepoImpl struct {
	db *pg.DB
}

func NewGymImageRepoImpl(db *pg.DB) (*GymImageRepoImpl, error) {
	if db == nil {
		return nil, errors.New("*sql.DB dependency is nil")
	}

	return &GymImageRepoImpl{
		db: db,
	}, nil
}

func (r *GymImageRepoImpl) Create(ctx context.Context, gym entity.GymImages) errors.Response {
	reqID, _ := c.GetRequestIdFromContext(ctx)
	log.Info().Str("RequestID", reqID).Msg("create gym images started on repository layer")

	_, err := r.db.Model(&gym).Insert()
	if err != nil {
		log.Error().Str("RequestID", reqID).Msg(err.Error())
		return errors.NewResponseConflictError(err.Error())
	}
	return nil
}

func (r *GymImageRepoImpl) GetImagesId(ctx context.Context, gymId int) ([]entity.GymImages, errors.Response) {
	reqID, _ := c.GetRequestIdFromContext(ctx)
	log.Info().Str("RequestID", reqID).Msg("GetById started on repository layer")

	var gym []entity.GymImages
	err := r.db.Model(&gym).Where("gym_id = ?", gymId).Select()
	if err != nil {
		log.Error().Str("RequestID", reqID).Msg(err.Error())
		return nil, errors.NewResponseConflictError(err.Error())
	}
	return gym, nil
}
