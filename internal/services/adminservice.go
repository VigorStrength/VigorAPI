package services

import (
	"context"
	"errors"

	"github.com/GhostDrew11/vigor-api/internal/db"
	"github.com/GhostDrew11/vigor-api/internal/models"
	"github.com/GhostDrew11/vigor-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminService struct {
	collection db.MongoCollection
	hasher utils.HashPasswordService
}

func NewAdminService(collection db.MongoCollection, hasher utils.HashPasswordService) *AdminService {
	return &AdminService{collection: collection, hasher: hasher}
}

func (as *AdminService) RegisterAdmin(ctx context.Context, admin models.Admin) error {
	// check if the admin already exists
	filter := bson.M{"email": admin.Email}
	count, err := as.collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("admin already exists")
	}

	// Hash the admin's password
	hashedPassword, err := as.hasher.HashPassword(admin.PasswordHash)
	if err != nil {
		return err
	}
	admin.PasswordHash = hashedPassword

	// Set the user ID
	admin.ID = primitive.NewObjectID()

	// Insert the admin into the database
	_, err = as.collection.InsertOne(ctx, admin)
	if err != nil {
		return err
	}

	return nil
}

func (as *AdminService) GetAdminByEmail(ctx context.Context, email string) (*models.Admin, error) {
	var admin models.Admin
	filter := bson.M{"email": email}

	result := as.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return &models.Admin{}, result.Err()
	}

	err := result.Decode(&admin)
	if err != nil {
		return &models.Admin{}, err
	}

	return &admin, nil
}

