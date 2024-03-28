package services

import (
	"context"

	"github.com/GhostDrew11/vigor-api/internal/models"
)

type MongoUserService interface {
	RegisterUser(ctx context.Context, input models.UserRegistrationInput) error
	GetUserByEmail(ctx context.Context, email, password string) (*models.User, error)
}