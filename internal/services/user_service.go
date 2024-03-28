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
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserCredentials = errors.New("invalid email or password")
)

type UserService struct {
	collection db.MongoCollection
	hasher     utils.HashPasswordService
}

func NewUserService(collection db.MongoCollection, hasher utils.HashPasswordService) *UserService {
	return &UserService{collection: collection, hasher: hasher}
}

func (us *UserService) RegisterUser(ctx context.Context, input models.UserRegistrationInput) error {
	// check if the user already exists
	filter := bson.M{"email": input.Email}
	count, err := us.collection.CountDocuments(ctx, filter)
	if err != nil {
		return fmt.Errorf("error checking if user already exists: %w", err)
	}

	if count > 0 {
		return ErrUserAlreadyExists
	}

	// Create a new user
	user, err := models.NewUserfromInput(input)
	if err != nil {
		return fmt.Errorf("error creating user from input: %w", err)
	}
	
	// Insert the user into the database
	_, err = us.collection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("error inserting user into database: %w", err)
	}

	return nil
}

func (us *UserService) GetUserByEmail(ctx context.Context, email, password string) (*models.User, error) {
	var user models.User
	filter := bson.M{"email": email}

	result := us.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return &models.User{}, fmt.Errorf("error fetching user with email %s: %w", email, result.Err())
	}

	err := result.Decode(&user)
	if err != nil {
		return &models.User{}, fmt.Errorf("error decoding user data for email %s: %w", email, err)
	}

	if !us.hasher.CheckPasswordHash(password, user.PasswordHash) {
		return &models.User{}, ErrInvalidUserCredentials
	}	

	return &user, nil
}