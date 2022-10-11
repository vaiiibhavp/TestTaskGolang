//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/scalent-io/healthapi/internal/middleware"

	gopg "github.com/scalent-io/healthapi/pkg/db/gopg"
	"github.com/scalent-io/healthapi/services/gym/repo"
	"github.com/scalent-io/healthapi/services/gym/service"
	"github.com/scalent-io/healthapi/services/gym/web"
)

var GymModuleSet = wire.NewSet(
	NewServerConfig,
	NewMiddlewareConfig,
	NewDBConfig,
	NewImageConfig,

	repo.NewGymRepoImpl,
	wire.Bind(new(service.GymRepo), new(*repo.GymRepoImpl)),

	repo.NewGymImageRepoImpl,
	wire.Bind(new(service.GymImagesRepo), new(*repo.GymImageRepoImpl)),

	service.NewGymServiceImpl,
	wire.Bind(new(service.GymService), new(*service.GymServiceImpl)),

	middleware.NewMiddlewareImpl,
	wire.Bind(new(middleware.Middleware), new(*middleware.MiddlewareImpl)),

	gopg.NewSqlDB,
	wire.Struct(new(web.GymHandlerRegistryOptions), "*"),
	web.NewHandlerRegistry,
)

func initServer() (*web.GymHandlerRegistry, error) {
	panic(wire.Build(GymModuleSet))
}
