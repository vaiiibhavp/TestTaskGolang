package utils

import (
	"encoding/json"
	"net/http"

	"github.com/scalent-io/healthapi/apimodel"
	"github.com/scalent-io/healthapi/pkg/errors"
)

const (
	STATUS_SUCCESS               = "success"
	STATUS_FAILED                = "failed"
	STATUS_INTERNAL_SERVER_ERROR = "internal server error"
)

func SendResponseWithData(w http.ResponseWriter, statusCode int, msg string, payload interface{}) {
	res := apimodel.Response{
		StatusCode: statusCode,
		Status:     STATUS_SUCCESS,
		Message:    msg,
		Data:       payload,
	}
	response, _ := json.Marshal(res)
	w.WriteHeader(statusCode)
	w.Write(response)
}

func SendResponseWithError(w http.ResponseWriter, err errors.Response, payload interface{}) {

	if err == nil {
		SendResponseWithError(w, errors.NewResponseInternalServerError(STATUS_INTERNAL_SERVER_ERROR), nil)
		return
	}
	res := apimodel.Response{
		StatusCode: err.StatusCode(),
		Status:     STATUS_FAILED,
		Message:    err.Error(),
		Data:       payload,
	}
	response, _ := json.Marshal(res)
	w.WriteHeader(err.StatusCode())
	w.Write(response)
}
