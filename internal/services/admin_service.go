package services

import (
	"context"
	"errors"
	"fmt"

	"github.com/GhostDrew11/vigor-api/internal/db"
	"github.com/GhostDrew11/vigor-api/internal/models"
	"github.com/GhostDrew11/vigor-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	ErrAdminAlreadyExists = errors.New("admin already exists")
	ErrInvalidAdminCredentials = errors.New("invalid email or password")
)

type AdminService struct {
	database db.MongoDatabase
	hasher utils.HashPasswordService
	parser utils.ParserService
}

func NewAdminService(database db.MongoDatabase, hasher utils.HashPasswordService, parser utils.ParserService) *AdminService {
	return &AdminService{database: database, hasher: hasher, parser: parser}
}

func (as *AdminService) RegisterAdmin(ctx context.Context, input models.AdminRegistrationInput) error {
	//Get the admin collection
	adminCollection := as.database.Collection("admins")

	// check if the admin already exists
	filter := bson.M{"email": input.Email}
	count, err := adminCollection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("error checking if admin already exists: %w", err)
	}

	if count > 0 {
		return ErrAdminAlreadyExists
	}

	// Create a new admin
	admin, err := models.NewAdminfromInput(input)
	if err != nil {
		return fmt.Errorf("error creating admin from input: %w", err)
	}

	// Insert the admin into the database
	_, err = adminCollection.InsertOne(ctx, admin)
	if err != nil {
		return fmt.Errorf("error inserting admin into database: %w", err)
	}

	return nil
}

func (as *AdminService) GetAdminByEmail(ctx context.Context, email, password string) (*models.Admin, error) {
	//Get the admin collection
	adminCollection := as.database.Collection("admins")

	var admin models.Admin
	filter := bson.M{"email": email}

	result := adminCollection.FindOne(ctx, filter)
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

