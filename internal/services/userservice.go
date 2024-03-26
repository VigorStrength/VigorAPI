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

type UserService struct {
	collection db.MongoCollection
	hasher     utils.HashPasswordService
}

func NewUserService(collection db.MongoCollection, hasher utils.HashPasswordService) *UserService {
	return &UserService{collection: collection, hasher: hasher}
}

func (us *UserService) RegisterUser(ctx context.Context, user models.User) error {
	// check if the user already exists
	filter := bson.M{"email": user.Email}
	count, err := us.collection.CountDocuments(ctx, filter)
	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("user already exists")
	}

	// Hash the user's password
	hashedPassword, err := us.hasher.HashPassword(user.PasswordHash)
	if err != nil {
		return err
	}
	user.PasswordHash = hashedPassword

	// Set the user ID
	user.ID = primitive.NewObjectID()

	// Insert the user into the database
	_, err = us.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	filter := bson.M{"email": email}

	result := us.collection.FindOne(ctx, filter)
	if result.Err() != nil {
		return &models.User{}, result.Err()
	}

	err := result.Decode(&user)
	if err != nil {
		return &models.User{}, err
	}

	return &user, nil
}