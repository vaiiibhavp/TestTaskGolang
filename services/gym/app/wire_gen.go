// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/scalent-io/healthapi/internal/middleware"
	"github.com/scalent-io/healthapi/pkg/db/gopg"
	"github.com/scalent-io/healthapi/services/gym/repo"
	"github.com/scalent-io/healthapi/services/gym/service"
	"github.com/scalent-io/healthapi/services/gym/web"
)

// Injectors from wire.go:

func initServer() (*web.GymHandlerRegistry, error) {
	dbConfig := NewDBConfig()
	db, err := gopg.NewSqlDB(dbConfig)
	if err != nil {
		return nil, err
	}
	gymRepoImpl, err := repo.NewGymRepoImpl(db)
	if err != nil {
		return nil, err
	}
	gymImageRepoImpl, err := repo.NewGymImageRepoImpl(db)
	if err != nil {
		return nil, err
	}
	gymServiceImpl, err := service.NewGymServiceImpl(gymRepoImpl, gymImageRepoImpl)
	if err != nil {
		return nil, err
	}
	middlewareConfig := NewMiddlewareConfig()
	middlewareImpl, err := middleware.NewMiddlewareImpl(middlewareConfig)
	if err != nil {
		return nil, err
	}
	serverConfig := NewServerConfig()
	gymHandlerRegistryOptions := web.GymHandlerRegistryOptions{
		GymService: gymServiceImpl,
		Middleware: middlewareImpl,
		Config:     serverConfig,
	}
	gymHandlerRegistry := web.NewHandlerRegistry(gymHandlerRegistryOptions)
	return gymHandlerRegistry, nil
}

// wire.go:

var GymModuleSet = wire.NewSet(
	NewServerConfig,
	NewMiddlewareConfig,
	NewDBConfig, repo.NewGymRepoImpl, wire.Bind(new(service.GymRepo), new(*repo.GymRepoImpl)), repo.NewGymImageRepoImpl, wire.Bind(new(service.GymImagesRepo), new(*repo.GymImageRepoImpl)), service.NewGymServiceImpl, wire.Bind(new(service.GymService), new(*service.GymServiceImpl)), middleware.NewMiddlewareImpl, wire.Bind(new(middleware.Middleware), new(*middleware.MiddlewareImpl)), gopg.NewSqlDB, wire.Struct(new(web.GymHandlerRegistryOptions), "*"), web.NewHandlerRegistry,
)
