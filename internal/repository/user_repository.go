package repository

import (
	"context"
	"database/sql"

	"github.com/Fox1N69/iq-testtask/internal/domain/entity"
	"github.com/Fox1N69/iq-testtask/pkg/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type UserRepository interface {
	UserByID(ctx context.Context, userID int64) (*entity.User, error)
}

type userRepository struct {
	db  *pgxpool.Pool
	log logger.Logger
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db:  db,
		log: logger.GetLogger(),
	}
}

func (r *userRepository) UserByID(ctx context.Context, userID int64) (*entity.User, error) {
	const query = `SELECT id, balance FROM users WHERE id = $1`

	var user entity.User

	row := r.db.QueryRow(ctx, query, userID)
	err := row.Scan(&user.ID, &user.Balance)
	if err != nil {
		if err == sql.ErrNoRows {
			r.log.WithField("userID", userID).Error("User not found")
			return nil, err
		}
		r.log.WithFields(logrus.Fields{
			"userID": userID,
			"error":  err,
		}).Error("Failed to scan user data")
		return nil, err
	}

	return &user, nil
}
