package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/cors"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	contextConstants "github.com/scalent-io/healthapi/pkg/context"
	contextModel "github.com/scalent-io/healthapi/pkg/context/models"
)

const (
	MSG_AUTH_MIDDLEWARE_STARTED = "auth middleware started"
	SESSION_ENDPOINT            = "/session"
)

type MiddlewareConfig struct {
	AuthServiceEndpoint string
}

// Note. Auth middleware needs to be implemented
type Middleware interface {
	Request() func(next http.Handler) http.Handler
	Cors() func(next http.Handler) http.Handler
}

type MiddlewareImpl struct {
	middlewareConfig *MiddlewareConfig
}

func NewMiddlewareImpl(middlewareConfig *MiddlewareConfig) (*MiddlewareImpl, error) {

	if middlewareConfig == nil {
		return nil, errors.New("middleware config  dependency is passed with nil value")
	}

	return &MiddlewareImpl{
		middlewareConfig: middlewareConfig,
	}, nil

}

func (m *MiddlewareImpl) Request() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			// get the reqId from the header
			reqId := r.Header.Get(contextConstants.REQUEST_ID)
			if reqId == EMPTY_STRING || len(reqId) != 36 {
				log.Debug().Str("RequestID", EMPTY_STRING).Msg(MSG_REQ_ID_IS_EMPTY)
				// generate request id
				newRequestId := uuid.New().String()
				// create context with values
				ctx = context.WithValue(ctx, contextModel.ContextKey(contextConstants.REQUEST_ID), newRequestId)
				// pass on the context to next request
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}
			// create context with values
			ctx = context.WithValue(ctx, contextModel.ContextKey(contextConstants.REQUEST_ID), reqId)
			// pass on the context to next request
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}

func (m *MiddlewareImpl) Cors() func(next http.Handler) http.Handler {

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Requested-With", "accept-version", "token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	return c.Handler

}
