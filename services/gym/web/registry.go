package web

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"github.com/scalent-io/healthapi/internal/middleware"

	"github.com/scalent-io/healthapi/pkg/server"
	"github.com/scalent-io/healthapi/services/gym/service"
)

type GymHandlerRegistryOptions struct {
	GymService service.GymService

	Middleware  middleware.Middleware
	Config      *server.Config
	ImageConfig *server.ImageConfig
}

type GymHandlerRegistry struct {
	options GymHandlerRegistryOptions
}

func NewHandlerRegistry(options GymHandlerRegistryOptions) *GymHandlerRegistry {
	return &GymHandlerRegistry{
		options: options,
	}
}

func fileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

func (h GymHandlerRegistry) StartServer() error {
	// register routes
	router, _ := h.RegisterRoutesTo()

	// fs := http.FileServer(http.Dir("uploads"))
	// router.Handle("/static/*", http.StripPrefix("/static/", fs))

	workDir, _ := os.Getwd()
	public := http.Dir(filepath.Join(workDir, "./", "uploads"))
	fileServer(router, "/static", public) // we can access assets folder with the name of static

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

	// r.Use(chimiddleware.Recoverer)
	r.Use(h.options.Middleware.Cors())
	r.Use(h.options.Middleware.Request())

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/gym", CreateGymHandler(h.options.GymService))
		r.Get("/gym", GetAllGymHandler(h.options.GymService, h.options.ImageConfig))
		r.Get("/gym/{id}", GetGymByIdHandler(h.options.GymService, h.options.ImageConfig))
		r.Get("/gym/search", SearchGymHandler(h.options.GymService, h.options.ImageConfig))
		r.Post("/gym/{id}", UploadGymImagesHandler(h.options.GymService, h.options.ImageConfig))
		r.Post("/gym/{id}/logo", UploadLogoHandler(h.options.GymService, h.options.ImageConfig))
	})
	return r, nil
}
