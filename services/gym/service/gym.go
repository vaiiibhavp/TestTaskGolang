package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/scalent-io/healthapi/apimodel"
	gymapimodel "github.com/scalent-io/healthapi/apimodel/gym"
	"github.com/scalent-io/healthapi/entity"
	c "github.com/scalent-io/healthapi/pkg/context"
	"github.com/scalent-io/healthapi/pkg/errors"
	"github.com/scalent-io/healthapi/pkg/server"
)

type GymServiceImpl struct {
	gymRepo       GymRepo
	gymImagesRepo GymImagesRepo
}

// constructor func
func NewGymServiceImpl(gymRepo GymRepo, gymImagesRepo GymImagesRepo) (*GymServiceImpl, error) {
	if gymRepo == nil {
		return nil, errors.New("gymRepo dependency is nil")
	} else if gymImagesRepo == nil {
		return nil, errors.New("gymImagesRepo dependency is nil")
	}
	return &GymServiceImpl{
		gymRepo:       gymRepo,
		gymImagesRepo: gymImagesRepo,
	}, nil
}

// Create add a new gym
func (s *GymServiceImpl) Create(ctx context.Context, gym *entity.Gym) (*apimodel.Response, errors.Response) {
	reqID, _ := c.GetRequestIdFromContext(ctx)
	log.Info().Str("RequestID", reqID).Msg("CreateGym started on service layer")

	gym.CreatedOn = time.Now()
	gym.ModifiedOn = time.Now()

	err := s.gymRepo.Create(ctx, gym)
	if err != nil {
		log.Error().Str("RequestID", reqID).Msg(err.Error())
		return nil, err
	}
	response := &apimodel.Response{
		StatusCode: http.StatusOK,
		Status:     STATUS_SUCCESS,
		Message:    MSG_GYM_CREATED,
		Data:       nil,
	}
	return response, nil
}

// GetAll return all the gyms
func (s *GymServiceImpl) GetAll(ctx context.Context, params url.Values, imageConfig *server.ImageConfig) (*gymapimodel.GetAllGymResponse, errors.Response) {
	reqID, _ := c.GetRequestIdFromContext(ctx)
	log.Info().Str("RequestID", reqID).Msg("GetAll started on service layer")

	page := params.Get("page")
	limit := params.Get("limit")

	limitInt, err := strconv.Atoi(limit)
	if err != nil {
		log.Error().Str("RequestID", reqID).Msg(err.Error())
		return nil, errors.NewResponseError(http.StatusBadRequest, "should have vaild limit")
	}
	pageInt, err := strconv.Atoi(page)
	if err != nil {
		log.Error().Str("RequestID", reqID).Msg(err.Error())
		return nil, errors.NewResponseError(http.StatusBadRequest, "should have vaild page no")
	}
	// Set default pagination values

	if limitInt == 0 && pageInt == 0 {
		limitInt = 20
		pageInt = 1
	}

	// calculate offset
	offset := CalcOffset(pageInt, limitInt)

	gyms, getGymsErr := s.gymRepo.GetGymImages(ctx, limitInt, offset)
	if getGymsErr != nil {
		log.Error().Str("RequestID", reqID).Msg(getGymsErr.Error())
		return nil, getGymsErr
	}

	//totalCount, countErr := s.gymRepo.Count(ctx)
	totalCount, countErr := s.gymRepo.CountGetGymImages(ctx)
	if countErr != nil {
		log.Error().Str("RequestID", reqID).Msg(countErr.Error())
		return nil, countErr
	}
	response := &gymapimodel.GetAllGymResponse{
		Response: apimodel.Response{
			StatusCode: http.StatusOK,
			Status:     STATUS_SUCCESS,
			Message:    GET_ALL_GYM_SUCCESS,
			Data:       gyms,
		},
		ImageBaseURL: imageConfig.ImageBaseURL,
		TotalRecords: totalCount,
		Page:         pageInt,
		Limit:        limitInt,
	}
	return response, nil
}

// GetById return single gym
func (s *GymServiceImpl) GetById(ctx context.Context, gymID int, imageConfig *server.ImageConfig) (*apimodel.Response, errors.Response) {
	reqID, _ := c.GetRequestIdFromContext(ctx)
	log.Info().Str("RequestID", reqID).Msg("GetById started on service layer")

	gym, getGymsErr := s.gymRepo.GetById(ctx, gymID)
	if getGymsErr != nil && !strings.Contains(getGymsErr.Error(), "no rows in result set") {
		log.Error().Str("RequestID", reqID).Msg(getGymsErr.Error())
		return nil, getGymsErr
	}
	if gym == nil {
		return &apimodel.Response{
			StatusCode: http.StatusOK,
			Status:     STATUS_FAILED,
			Message:    GET_GYM_BY_ID_FAILED,
		}, nil
	}

	gymImages, getGymImagesErr := s.gymImagesRepo.GetImagesId(ctx, gymID)
	if getGymImagesErr != nil && !strings.Contains(getGymImagesErr.Error(), "no rows in result set") {
		log.Error().Str("RequestID", reqID).Msg(getGymImagesErr.Error())
		return nil, getGymImagesErr
	}

	data := gymapimodel.GymDetailsResponse{
		ImageBaseURL: imageConfig.ImageBaseURL,
		Gym:          gym,
		GymImages:    gymImages,
	}

	response := &apimodel.Response{
		StatusCode: http.StatusOK,
		Status:     STATUS_SUCCESS,
		Message:    GET_GYM_BY_ID_SUCCESS,
		Data:       data,
	}
	return response, nil
}

func CalcOffset(pageNum int, limit int) int {
	return (pageNum - 1) * limit
}

// Search find nearest gyms based on long and lat

func (s *GymServiceImpl) Search(ctx context.Context, params url.Values, imageConfig *server.ImageConfig) (*apimodel.Response, errors.Response) {

	reqID, _ := c.GetRequestIdFromContext(ctx)
	log.Info().Str("RequestID", reqID).Msg("Search started on service layer")

	fmt.Println("params.Get(lat)", params.Get("lat"))

	lat, err := strconv.ParseFloat(params.Get("lat"), 64)
	if err != nil {
		log.Error().Str("RequestID", reqID).Msg(err.Error())
		//return nil, err
	}

	long, err := strconv.ParseFloat(params.Get("lng"), 64)
	if err != nil {
		log.Error().Str("RequestID", reqID).Msg(err.Error())
	}

	gyms, getGymsErr := s.gymRepo.Search(ctx, lat, long)
	if getGymsErr != nil {
		log.Error().Str("RequestID", reqID).Msg(getGymsErr.Error())
		return nil, getGymsErr
	}

	searchResponse := gymapimodel.GymSearchResponse{
		ImageBaseURL: imageConfig.ImageBaseURL,
		GymDistance:  gyms,
	}

	response := &apimodel.Response{
		StatusCode: http.StatusOK,
		Status:     STATUS_SUCCESS,
		Message:    GYM_SEARCH_SUCCESS,
		Data:       searchResponse,
	}

	return response, nil
}

// Upload

func (s *GymServiceImpl) Upload(ctx context.Context, images []entity.GymImages) (*apimodel.Response, errors.Response) {
	reqID, _ := c.GetRequestIdFromContext(ctx)
	log.Info().Str("RequestID", reqID).Msg("Upload started on service layer")

	for _, img := range images {
		img.CreatedOn = time.Now()

		err := s.gymImagesRepo.Create(ctx, img)
		if err != nil {
			log.Error().Str("RequestID", reqID).Msg(err.Error())
			return nil, err
		}
	}
	response := &apimodel.Response{
		StatusCode: http.StatusOK,
		Status:     STATUS_SUCCESS,
		Message:    MSG_GYM_IMAGE_UPLOD,
	}
	return response, nil
}

// Upload Logo

func (s *GymServiceImpl) UploadLogo(ctx context.Context, req gymapimodel.CreateLogoReq) (*apimodel.Response, errors.Response) {
	reqID, _ := c.GetRequestIdFromContext(ctx)
	log.Info().Str("RequestID", reqID).Msg("UploadLogo started on service layer")

	resp, getErr := s.gymRepo.GetById(ctx, req.GymID)
	if getErr != nil {
		log.Error().Str("RequestID", reqID).Msg(getErr.Error())
		return nil, getErr
	}
	resp.LogoName = req.LogoName
	resp.ModifiedOn = time.Now()

	updateErr := s.gymRepo.Update(ctx, resp)
	if updateErr != nil {
		log.Error().Str("RequestID", reqID).Msg(updateErr.Error())
		return nil, updateErr
	}

	response := &apimodel.Response{
		StatusCode: http.StatusOK,
		Status:     STATUS_SUCCESS,
		Message:    MSG_GYM_IMAGE_UPLOD,
	}
	return response, nil
}
