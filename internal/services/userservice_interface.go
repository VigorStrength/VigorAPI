package services

import (
	"context"

	"github.com/GhostDrew11/vigor-api/internal/models"
)

type MongoUserService interface {
	RegisterUser(ctx context.Context, user models.User) error
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}