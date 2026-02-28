//go:build wireinject
// +build wireinject

package app

import (
	"github.com/farid141/go-rest-api/config"
	"github.com/farid141/go-rest-api/controller"
	"github.com/farid141/go-rest-api/router"
	"github.com/sirupsen/logrus"

	"github.com/farid141/go-rest-api/db"
	"github.com/farid141/go-rest-api/repository"
	"github.com/farid141/go-rest-api/service"
	"github.com/google/wire"
)

type AppContainer struct {
	UserService service.UserService
	AuthService service.AuthService
	PostService service.PostService
	Router      *router.Router
	Logger      *logrus.Logger
	Config      *config.Config
}

func InitializeApp() (*AppContainer, error) {
	wire.Build(
		config.ProviderSet,
		db.ProviderSet,
		repository.ProviderSet,
		service.ProviderSet,
		controller.ProviderSet,
		router.ProviderSet,
		wire.Struct(new(AppContainer), "*"),
	)
	return &AppContainer{}, nil
}
