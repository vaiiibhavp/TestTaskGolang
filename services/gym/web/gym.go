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
	"github.com/scalent-io/healthapi/pkg/server"
	"github.com/scalent-io/healthapi/pkg/utils"
	"github.com/scalent-io/healthapi/pkg/validation"
	"github.com/scalent-io/healthapi/services/gym/service"
)

const MAX_UPLOAD_SIZE = (1024 * 1024) * 5

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

func GetAllGymHandler(gymService service.GymService, imageConfig *server.ImageConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID, _ := c.GetRequestIdFromContext(r.Context())
		log.Info().Str("RequestID", reqID).Msg("GetAllGym handler")

		params := r.URL.Query()

		gyms, getGymsErr := gymService.GetAll(r.Context(), params, imageConfig)
		if getGymsErr != nil {
			log.Error().Str("RequestID", reqID).Msg(getGymsErr.Error())
			utils.SendResponseWithError(w, getGymsErr, nil)
			return
		}
		render.JSON(w, r, gyms)
	}
}

func GetGymByIdHandler(gymService service.GymService, imageConfig *server.ImageConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID, _ := c.GetRequestIdFromContext(r.Context())
		log.Info().Str("RequestID", reqID).Msg("GetGymById handler")

		id := chi.URLParam(r, "id")
		gymId, err := strconv.Atoi(id)
		if err != nil {
			log.Error().Str("RequestID", reqID).Msg(err.Error())
			return
		}

		gym, getGymErr := gymService.GetById(r.Context(), gymId, imageConfig)
		if getGymErr != nil {
			log.Error().Str("RequestID", reqID).Msg(getGymErr.Error())
			http.Error(w, "gymId not found", http.StatusBadRequest)
			return
		}

		if gym == nil {
			http.Error(w, "gymId not found", http.StatusBadRequest)
			return
		}
		render.JSON(w, r, gym)
	}
}

func SearchGymHandler(gymService service.GymService, imageConfig *server.ImageConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID, _ := c.GetRequestIdFromContext(r.Context())
		log.Info().Str("RequestID", reqID).Msg("SearchGym Handler")
		// get lat and lng as query params
		params := r.URL.Query()

		fmt.Println("params", params)

		data, searchErr := gymService.Search(r.Context(), params, imageConfig)
		if searchErr != nil {
			log.Error().Str("RequestID", reqID).Msg(searchErr.Error())
			utils.SendResponseWithError(w, searchErr, nil)
			return
		}
		render.JSON(w, r, data)
	}
}

func UploadGymImagesHandler(gymService service.GymService, imageConfig *server.ImageConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		reqID, _ := c.GetRequestIdFromContext(r.Context())
		log.Info().Str("RequestID", reqID).Msg("Upload GymImages Handler")

		id := chi.URLParam(r, "id")
		gymId, err := strconv.Atoi(id)
		if err != nil {
			log.Error().Str("RequestID", reqID).Msg(err.Error())
			return
		}

		params := r.URL.Query()
		var imageType string
		val, _ := strconv.Atoi(params.Get("type")) // @Todo: Need to handle error
		if val == 2 {
			imageType = "amenities"
		}

		//@Todo: The logic should be moved to service layer
		res, getError := gymService.GetById(r.Context(), gymId, imageConfig)
		if getError != nil {
			log.Error().Str("RequestID", reqID).Msg(getError.Error())
			http.Error(w, "gymId not found", http.StatusBadRequest)
			return
		}

		// return appropriate message if not exist
		if res == nil {
			http.Error(w, "gymId not found", http.StatusBadRequest)
			return
		}

		var gymImageReq []apimodel.CreateGymImageReq

		//For time being file upload logic kept in web layer
		//@Todo: The logic should be moved to repo layer
		//@Todo: Instead storing images to local folder we have to upload to AWS S3 buckets
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
			fname := fmt.Sprintf("%s%s", uuid, filepath.Ext(fileHeader.Filename))
			fileName := fmt.Sprintf("./uploads/%s", fname)
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
				Label:     fname,
				Type:      imageType,
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
		log.Info().Str("RequestID", reqID).Msg("UploadGymImages completed at handler level")
		render.JSON(w, r, data)
	}
}

func UploadLogoHandler(gymService service.GymService, imageConfig *server.ImageConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqID, _ := c.GetRequestIdFromContext(r.Context())
		log.Info().Str("RequestID", reqID).Msg("Upload GymImages Handler")

		id := chi.URLParam(r, "id")
		gymId, err := strconv.Atoi(id)
		if err != nil {
			log.Error().Str("RequestID", reqID).Msg(err.Error())
			return
		}

		//@Todo: The logic should be moved to service layer
		res, getError := gymService.GetById(r.Context(), gymId, imageConfig)
		if getError != nil {
			log.Error().Str("RequestID", reqID).Msg(getError.Error())
			http.Error(w, "gymId not found", http.StatusBadRequest)
			return
		}

		if res.Status == "Failed" {
			render.JSON(w, r, res)
			return
		}

		//For time being file upload logic kept in web layer
		//@Todo: The logic should be moved to repo layer
		//@Todo: Instead storing images to local folder we have to upload to AWS S3 buckets
		r.ParseMultipartForm(32 << 20)
		file, handler, err := r.FormFile("logo")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer file.Close()

		// For time being upload file path kept as hard coded later will take it from config.yml
		uuid := uuid.New().String()
		fname := fmt.Sprintf("logos/%s%s", uuid, filepath.Ext(handler.Filename))
		f, err := os.OpenFile("./uploads/"+fname, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		logoRequest := apimodel.CreateLogoReq{
			GymID:    gymId,
			LogoName: fname,
		}

		data, err := gymService.UploadLogo(r.Context(), logoRequest)
		if err != nil {
			log.Error().Str("RequestID", reqID).Msg(err.Error())
			utils.SendResponseWithError(w, errors.NewResponseInternalServerError("internal server error"), nil)
			return
		}
		log.Info().Str("RequestID", reqID).Msg("UploadLogo completed at handler level")
		render.JSON(w, r, data)
	}
}
