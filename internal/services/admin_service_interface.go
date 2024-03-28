package services

import (
	"context"

	"github.com/GhostDrew11/vigor-api/internal/models"
)

type MongoAdminService interface {
	RegisterAdmin(ctx context.Context, input models.AdminRegistrationInput) error
	GetAdminByEmail(ctx context.Context, email, password string) (*models.Admin, error)
}