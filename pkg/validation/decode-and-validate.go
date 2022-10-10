package validation

import (
	"encoding/json"
	"io"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/schema"
	"github.com/scalent-io/healthapi/pkg/errors"
)

const (
	MSG_INVALID_REQUEST = "invalid request"
)

var decoder = schema.NewDecoder()

func DecodeAndVaildate(r io.Reader, requestInstance interface{}) errors.Response {
	// decode the request
	err := json.NewDecoder(r).Decode(requestInstance)
	if err != nil {
		return errors.NewResponseBadRequestError(err.Error())
	}

	// validate the request
	validate := validator.New()
	err = validate.Struct(requestInstance)
	if err != nil {
		return errors.NewResponseBadRequestError(err.Error())
	}

	return nil
}

func DecodeAndVaildateQueryparams(r map[string][]string, requestInstance interface{}) errors.Response {

	err := decoder.Decode(requestInstance, r)
	if err != nil {
		return errors.NewResponseBadRequestError(err.Error())
	}

	// validate the request
	validate := validator.New()
	err = validate.Struct(requestInstance)
	if err != nil {
		return errors.NewResponseBadRequestError(err.Error())
	}

	return nil
}
