package handler

import (
	"github.com/Fox1N69/iq-testtask/internal/domain/service"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserHandler,
	NewTransactionHandler,
	service.ProviderSet,
)
