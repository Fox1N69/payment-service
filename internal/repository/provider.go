package repository

import (
	"github.com/Fox1N69/iq-testtask/storage/postgres"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(
	postgres.ProviderSet,
	NewUserRepository,
	NewTransactionRepository,
)
