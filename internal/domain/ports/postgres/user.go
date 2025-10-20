package postgres

import "hona/backend/internal/domain/entities"

type UserRepository interface {
	FindUserByEmail(email string) (*entities.User, error)
}
