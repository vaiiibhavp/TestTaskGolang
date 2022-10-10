package web

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog/log"
	"github.com/scalent-io/healthapi/internal/middleware"

	"github.com/scalent-io/healthapi/pkg/server"
	"github.com/scalent-io/healthapi/services/gym/service"
)

type GymHandlerRegistryOptions struct {
	GymService service.GymService

	Middleware middleware.Middleware
	Config     *server.Config
}

type GymHandlerRegistry struct {
	options GymHandlerRegistryOptions
}

func NewHandlerRegistry(options GymHandlerRegistryOptions) *GymHandlerRegistry {
	return &GymHandlerRegistry{
		options: options,
	}
}

func (h GymHandlerRegistry) StartServer() error {
	// register routes
	router, _ := h.RegisterRoutesTo()

	//Start Http Server
	log.Printf("server running on port %d", h.options.Config.Port)
	httpServer := server.New(router, server.Port(fmt.Sprintf("%d", h.options.Config.Port)))

	// initialize log caller
	log.Logger = log.With().Caller().Logger()

	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		errMessage := fmt.Errorf("app - Run - : %v", s)
		log.Error().Msg(errMessage.Error())
	case err := <-httpServer.Notify():
		errMessage := fmt.Errorf("app - Run - httpServer.Notify: %w", err)
		log.Error().Msg(errMessage.Error())
	}

	// Shutdown
	err := httpServer.Shutdown()
	if err != nil {
		errMessage := fmt.Errorf("app - Run - httpServer.Shutdown: %w", err)
		log.Error().Msg(errMessage.Error())
		return errMessage
	}
	return nil
}

func (h GymHandlerRegistry) RegisterRoutesTo() (*chi.Mux, error) {
	r := chi.NewRouter()

	r.Use(chimiddleware.Recoverer)
	r.Use(h.options.Middleware.Cors())

	r.Use(h.options.Middleware.Request())

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/gym", CreateGymHandler(h.options.GymService))
		r.Get("/gym", GetAllGymHandler(h.options.GymService))
		r.Get("/gym/{id}", GetGymByIdHandler(h.options.GymService))
		r.Get("/gym/search", SearchGymHandler(h.options.GymService))
		r.Post("/gym/{id}", UploadGymImagesHandler(h.options.GymService))
	})
	return r, nil
}
