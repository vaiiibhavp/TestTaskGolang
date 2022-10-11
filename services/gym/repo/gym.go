package repo

import (
	"context"

	"github.com/rs/zerolog/log"

	"github.com/go-pg/pg/v10"

	"github.com/scalent-io/healthapi/entity"
	c "github.com/scalent-io/healthapi/pkg/context"
	"github.com/scalent-io/healthapi/pkg/errors"
)

type GymRepoImpl struct {
	db *pg.DB
}

func NewGymRepoImpl(db *pg.DB) (*GymRepoImpl, error) {
	if db == nil {
		return nil, errors.New("*sql.DB dependency is nil")
	}

	return &GymRepoImpl{
		db: db,
	}, nil
}

func (r *GymRepoImpl) Create(ctx context.Context, gym *entity.Gym) errors.Response {
	reqID, _ := c.GetRequestIdFromContext(ctx)
	log.Info().Str("RequestID", reqID).Msg("create gym started on repository layer")

	_, err := r.db.Model(gym).Insert()
	if err != nil {
		log.Error().Str("RequestID", reqID).Msg(err.Error())
		return errors.NewResponseConflictError(err.Error())
	}
	return nil
}

func (r *GymRepoImpl) GetAll(ctx context.Context, limit, offset int) ([]entity.Gym, errors.Response) {
	reqID, _ := c.GetRequestIdFromContext(ctx)
	log.Info().Str("RequestID", reqID).Msg("GetAll started on repository layer")

	var gyms []entity.Gym
	err := r.db.Model(&gyms).Order("created_on DESC").Limit(limit).Offset(offset).Select()
	if err != nil {
		log.Error().Str("RequestID", reqID).Msg(err.Error())
		return nil, errors.NewResponseConflictError(err.Error())
	}
	return gyms, nil
}

func (r *GymRepoImpl) Count(ctx context.Context) (int, errors.Response) {
	reqID, _ := c.GetRequestIdFromContext(ctx)
	log.Info().Str("RequestID", reqID).Msg("GetAll started on repository layer")

	count, err := r.db.Model((*entity.Gym)(nil)).Count()
	if err != nil {
		log.Error().Str("RequestID", reqID).Msg(err.Error())
		return 0, errors.NewResponseConflictError(err.Error())
	}
	return count, nil

}

func (r *GymRepoImpl) GetById(ctx context.Context, gymId int) (*entity.Gym, errors.Response) {
	reqID, _ := c.GetRequestIdFromContext(ctx)
	log.Info().Str("RequestID", reqID).Msg("GetById started on repository layer")

	gym := new(entity.Gym)
	err := r.db.Model(gym).Where("id = ?", gymId).Select()
	if err != nil {
		log.Error().Str("RequestID", reqID).Msg(err.Error())
		return nil, errors.NewResponseConflictError(err.Error())
	}
	return gym, nil
}

func (r *GymRepoImpl) Search(ctx context.Context, lat, long float64) ([]entity.GymDistance, errors.Response) {
	reqID, _ := c.GetRequestIdFromContext(ctx)
	log.Info().Str("RequestID", reqID).Msg("Search started on repository layer")

	var gyms []entity.GymDistance

	query := `SELECT g.*, (point(g.long, g.lat) <@> point(?, ?)) * 1609.344 AS distance , gi.image_type, gi.label, gi.type
	FROM
	 gyms as g
	 LEFT JOIN gym_images as gi
	ON g.id = gi.gym_id 
	WHERE
		(point(g.long, g.lat) <@> point(?, ?)) < (10000 / 1609.344)
	ORDER BY
	 distance`

	_, err := r.db.Query(&gyms, query, long, lat, long, lat)
	if err != nil {
		log.Error().Str("RequestID", reqID).Msg(err.Error())
		return nil, errors.NewResponseConflictError(err.Error())
	}
	return gyms, nil
}

// func to get gym details with gym images

func (r *GymRepoImpl) GetGymImages(ctx context.Context, limit, offset int) ([]entity.GymDetails, errors.Response) {
	reqID, _ := c.GetRequestIdFromContext(ctx)
	log.Info().Str("RequestID", reqID).Msg("Search started on repository layer")

	var gyms []entity.GymDetails

	order := " order by g.created_on desc limit ? offset ?"

	query := `SELECT g.*, gi.image_type, gi.label, gi.type
	FROM gyms as g
	LEFT JOIN gym_images as gi
	ON g.id = gi.gym_id `

	query = query + order

	_, err := r.db.Query(&gyms, query, limit, offset)
	if err != nil {
		log.Error().Str("RequestID", reqID).Msg(err.Error())
		return nil, errors.NewResponseConflictError(err.Error())
	}
	return gyms, nil
}

// Count
func (r *GymRepoImpl) CountGetGymImages(ctx context.Context) (int, errors.Response) {
	reqID, _ := c.GetRequestIdFromContext(ctx)
	log.Info().Str("RequestID", reqID).Msg("Search started on repository layer")

	var count entity.GymCount

	query := `SELECT count(g.id) as total_records
	FROM gyms as g
	LEFT JOIN gym_images as gi
	ON g.id = gi.gym_id `

	_, err := r.db.Query(&count, query)
	if err != nil {
		log.Error().Str("RequestID", reqID).Msg(err.Error())
		return 0, errors.NewResponseConflictError(err.Error())
	}
	return count.TotalRecords, nil
}

func (r *GymRepoImpl) Update(ctx context.Context, gym *entity.Gym) errors.Response {
	reqID, _ := c.GetRequestIdFromContext(ctx)
	log.Info().Str("RequestID", reqID).Msg("create gym started on repository layer")

	_, err := r.db.Model(gym).WherePK().Update()
	if err != nil {
		log.Error().Str("RequestID", reqID).Msg(err.Error())
		return errors.NewResponseConflictError(err.Error())
	}
	return nil
}
