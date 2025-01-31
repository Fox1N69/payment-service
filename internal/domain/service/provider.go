package service

import (
	"github.com/Fox1N69/iq-testtask/internal/repository"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	NewUserService,
	NewTransactionService,
	repository.ProviderSet,
)
