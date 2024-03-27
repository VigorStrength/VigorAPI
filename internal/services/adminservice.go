package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/GhostDrew11/vigor-api/internal/db"
	"github.com/GhostDrew11/vigor-api/internal/models"
	"github.com/GhostDrew11/vigor-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	ErrAdminAlreadyExists = errors.New("admin already exists")
	ErrInvalidAdminCredentials = errors.New("invalid email or password")
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
		return fmt.Errorf("error checking if admin already exists: %w", err)
	}

	if count > 0 {
		return ErrAdminAlreadyExists
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

func (as *AdminService) GetAdminByEmail(ctx context.Context, email, password string) (*models.Admin, error) {
	var admin models.Admin
	filter := bson.M{"email": email}

	result := as.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return &models.Admin{}, fmt.Errorf("error fetching admin with email %s: %w", email, result.Err())
	}

	err := result.Decode(&admin)
	if err != nil {
		return &models.Admin{}, fmt.Errorf("error decoding admin data for email %s: %w", email, err)
	}

	if !as.hasher.CheckPasswordHash(password, admin.PasswordHash) {
		return &models.Admin{}, ErrInvalidAdminCredentials
	}

	return &admin, nil
}

