package services

import (
	"context"

	"github.com/GhostDrew11/vigor-api/internal/models"
)

type MongoAdminService interface {
	RegisterAdmin(ctx context.Context, admin models.Admin) error
	GetAdminByEmail(ctx context.Context, email string) (*models.Admin, error)
}