//go:build wireinject
// +build wireinject

package server

import (
	"github.com/Fox1N69/iq-testtask/internal/config"
	"github.com/Fox1N69/iq-testtask/internal/delivery/http/handler"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	config.ProviderSet,
	handler.ProviderSet,
	New,
)

func InitializeServer() (*Server, error) {
	wire.Build(ProviderSet)
	return &Server{}, nil
}
