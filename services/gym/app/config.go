package main

import (
	"github.com/scalent-io/healthapi/internal/middleware"

	"github.com/scalent-io/healthapi/pkg/db/gopg"
	"github.com/scalent-io/healthapi/pkg/server"
	"github.com/spf13/viper"
)

var config struct {
	Server      server.Config
	Middleware  middleware.MiddlewareConfig
	DB          gopg.DbConfig
	ImageConfig server.ImageConfig
}

func initConfig() error {
	viper.AddConfigPath(".")
	// viper.AddConfigPath("./config/")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.SetEnvPrefix("auth")
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	viper.Unmarshal(&config)
	return nil
}

func NewServerConfig() *server.Config {
	return &config.Server
}

func NewMiddlewareConfig() *middleware.MiddlewareConfig {
	return &config.Middleware
}

func NewDBConfig() *gopg.DbConfig {
	return &config.DB
}

func NewImageConfig() *server.ImageConfig {
	return &config.ImageConfig
}
