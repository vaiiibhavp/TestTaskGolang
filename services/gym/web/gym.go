package web

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	apimodel "github.com/scalent-io/healthapi/apimodel/gym"
	"github.com/scalent-io/healthapi/internal/converter"
	c "github.com/scalent-io/healthapi/pkg/context"
	"github.com/scalent-io/healthapi/pkg/errors"
	"github.com/scalent-io/healthapi/pkg/utils"
	"github.com/scalent-io/healthapi/pkg/validation"
	"github.com/scalent-io/healthapi/services/gym/service"
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 5

func CreateGymHandler(gymService service.GymService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID, _ := c.GetRequestIdFromContext(r.Context())
		log.Info().Str("RequestID", reqID).Msg("CreateGym handler")

		createGymReq := &apimodel.CreateGymReq{}
		responseErr := validation.DecodeAndVaildate(r.Body, createGymReq)
		if responseErr != nil {
			log.Error().Str("RequestID", reqID).Msg(responseErr.Error())
			utils.SendResponseWithError(w, responseErr, nil)
			return
		}

		gym := converter.GymAPIToGymEntity(*createGymReq)

		data, err := gymService.Create(r.Context(), &gym)
		if err != nil {
			log.Error().Str("RequestID", reqID).Msg(err.Error())
			utils.SendResponseWithError(w, errors.NewResponseInternalServerError("internal server error"), nil)
			return
		}
		log.Info().Str("RequestID", reqID).Msg("CreateGym completed at handler level")
		render.JSON(w, r, data)
	}
}

func GetAllGymHandler(gymService service.GymService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID, _ := c.GetRequestIdFromContext(r.Context())
		log.Info().Str("RequestID", reqID).Msg("GetAllGym handler")

		params := r.URL.Query()

		gyms, getGymsErr := gymService.GetAll(r.Context(), params)
		if getGymsErr != nil {
			log.Error().Str("RequestID", reqID).Msg(getGymsErr.Error())
			utils.SendResponseWithError(w, getGymsErr, nil)
			return
		}
		render.JSON(w, r, gyms)
	}
}

func GetGymByIdHandler(gymService service.GymService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID, _ := c.GetRequestIdFromContext(r.Context())
		log.Info().Str("RequestID", reqID).Msg("GetGymById handler")

		id := chi.URLParam(r, "id")
		gymId, err := strconv.Atoi(id)
		if err != nil {
			log.Error().Str("RequestID", reqID).Msg(err.Error())
			return
		}

		gym, getGymErr := gymService.GetById(r.Context(), gymId)
		if getGymErr != nil {
			log.Error().Str("RequestID", reqID).Msg(getGymErr.Error())
			utils.SendResponseWithError(w, getGymErr, nil)
			return
		}
		render.JSON(w, r, gym)
	}
}

func SearchGymHandler(gymService service.GymService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID, _ := c.GetRequestIdFromContext(r.Context())
		log.Info().Str("RequestID", reqID).Msg("SearchGym Handler")
		// get lat and lng as query params
		params := r.URL.Query()

		fmt.Println("params", params)

		data, searchErr := gymService.Search(r.Context(), params)
		if searchErr != nil {
			log.Error().Str("RequestID", reqID).Msg(searchErr.Error())
			utils.SendResponseWithError(w, searchErr, nil)
			return
		}
		render.JSON(w, r, data)
	}
}

func UploadGymImagesHandler(gymService service.GymService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID, _ := c.GetRequestIdFromContext(r.Context())
		log.Info().Str("RequestID", reqID).Msg("Upload GymImages Handler")

		id := chi.URLParam(r, "id")
		gymId, err := strconv.Atoi(id)
		if err != nil {
			log.Error().Str("RequestID", reqID).Msg(err.Error())
			return
		}

		var gymImageReq []apimodel.CreateGymImageReq

		// 32 MB is the default used by FormFile()
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		files := r.MultipartForm.File["file"]

		for _, fileHeader := range files {
			if fileHeader.Size > MAX_UPLOAD_SIZE {
				http.Error(w, fmt.Sprintf("The uploaded image is too big: %s. Please use an image less than 5MB in size", fileHeader.Filename), http.StatusBadRequest)
				return
			}

			// Open the file
			file, err := fileHeader.Open()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			defer file.Close()

			buff := make([]byte, 512)
			_, err = file.Read(buff)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			filetype := http.DetectContentType(buff)
			if filetype != "image/jpeg" && filetype != "image/png" {
				http.Error(w, "The provided file format is not allowed. Please upload a JPEG or PNG image", http.StatusBadRequest)
				return
			}

			_, err = file.Seek(0, io.SeekStart)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// For time being upload file path kept as hard coded later will take it from config.yml
			err = os.MkdirAll("./uploads", os.ModePerm)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			uuid := uuid.New().String()

			fileName := fmt.Sprintf("./uploads/%s%s", uuid, filepath.Ext(fileHeader.Filename))
			f, err := os.Create(fileName)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			defer f.Close()

			_, err = io.Copy(f, file)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			images := apimodel.CreateGymImageReq{
				GymID:     gymId,
				ImageType: filetype,
				Label:     fileName,
			}
			gymImageReq = append(gymImageReq, images)
		}

		// call service to insert into database
		gymImages := converter.GymAPIToGymImagesEntity(gymImageReq)

		data, err := gymService.Upload(r.Context(), gymImages)
		if err != nil {
			log.Error().Str("RequestID", reqID).Msg(err.Error())
			utils.SendResponseWithError(w, errors.NewResponseInternalServerError("internal server error"), nil)
			return
		}
		log.Info().Str("RequestID", reqID).Msg("CreateGym completed at handler level")
		render.JSON(w, r, data)
	}
}
