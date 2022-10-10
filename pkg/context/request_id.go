package context

import (
	"context"
	"errors"

	"github.com/scalent-io/healthapi/pkg/context/models"
)

const (
	EMPTY_STRING                  = ""
	ERR_NO_VAL_REQUEST_ID         = "no value for the context key " + REQUEST_ID
	ERR_NO_VALID_VALUE_REQUEST_ID = "no valid value for the context key " + REQUEST_ID
	ERR_NO_VAL_SESSION            = "no value for the context key " + SESSION_DATA
	ERR_NO_VALID_VALUE_SESSION    = "no valid value for the context key " + SESSION_DATA
)

func GetRequestIdFromContext(ctx context.Context) (string, error) {
	reqID := ctx.Value(models.ContextKey(REQUEST_ID))
	if reqID == nil {
		return EMPTY_STRING, errors.New(ERR_NO_VAL_REQUEST_ID)
	}
	reqIDString, ok := reqID.(string)
	if !ok {
		return reqIDString, errors.New(ERR_NO_VALID_VALUE_REQUEST_ID)
	}
	return reqIDString, nil
}

// func GetSessionFromContext(ctx context.Context) (entity.Session, error) {
// 	session := ctx.Value(models.ContextKey(SESSION_DATA))
// 	if session == nil {
// 		return entity.Session{}, errors.New(ERR_NO_VAL_SESSION)
// 	}
// 	sessionVal, ok := session.(apimodel.Session)
// 	if !ok {
// 		return entity.Session{}, errors.New(ERR_NO_VALID_VALUE_SESSION)
// 	}
// 	return entity.Session{
// 		UserID:    sessionVal.UserID,
// 		Email:     sessionVal.Email,
// 		FirstName: sessionVal.FirstName,
// 		LastName:  sessionVal.LastName,
// 		Roles:     sessionVal.Roles,
// 		TTL:       sessionVal.TTL,
// 	}, nil
// }

// func GetTokenFromContext(ctx context.Context) (string, error) {
// 	token := ctx.Value(models.ContextKey(TOKEN))

// 	// fmt.Println("========= ")
// 	// fmt.Println(ctx.Value(models.ContextKey("token")))
// 	// fmt.Println("========= ")

// 	if token == nil {
// 		return "", errors.New("no valid value for the context key token")
// 	}

// 	tokenVal, ok := token.(string)
// 	if !ok {
// 		return "", errors.New("no valid value for the context key token")
// 	}
// 	return tokenVal, nil
// }
