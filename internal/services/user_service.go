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
	ErrUsernameAlreadyTaken = errors.New("username already taken")
	ErrInvalidUserCredentials = errors.New("invalid email or password")
)

type UserService struct {
	database db.MongoDatabase
	hasher utils.HashPasswordService
}

func NewUserService(database db.MongoDatabase, hasher utils.HashPasswordService) *UserService {
	return &UserService{database: database, hasher: hasher}
}

func (us *UserService) RegisterUser(ctx context.Context, input models.UserRegistrationInput) error {
	//Get the user collection
	userCollection := us.database.Collection("users")

	// check if the user already exists
	emailFilter := bson.M{"email": input.Email}
	emailCount, err := userCollection.CountDocuments(ctx, emailFilter)
	if err != nil {
		return fmt.Errorf("error checking if user already exists: %w", err)
	}

	if emailCount > 0 {
		return ErrUserAlreadyExists
	}

	// check if the username already exists
	usernameFilter := bson.M{"profileInformation.username": input.ProfileInformation.Username}
	usernameCount, err := userCollection.CountDocuments(ctx, usernameFilter)
	if err != nil {
		return fmt.Errorf("error checking if username already exists: %w", err)
	}

	if usernameCount > 0 {
		return ErrUsernameAlreadyTaken
	}

	// Create a new user
	user, err := models.NewUserfromInput(input)
	if err != nil {
		return fmt.Errorf("error creating user from input: %w", err)
	}

	// Insert the user into the database
	_, err = userCollection.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("error inserting user into database: %w", err)
	}

	return nil
}

func (us *UserService) GetUserByEmail(ctx context.Context, email, password string) (*models.User, error) {
	//Get the user collection
	userCollection := us.database.Collection("users")

	var user models.User
	filter := bson.M{"email": email}

	result := userCollection.FindOne(ctx, filter)
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

// type UserService struct {
// 	collection db.MongoCollection
// 	hasher     utils.HashPasswordService
// }

// func NewUserService(collection db.MongoCollection, hasher utils.HashPasswordService) *UserService {
// 	return &UserService{collection: collection, hasher: hasher}
// }

// func (us *UserService) RegisterUser(ctx context.Context, input models.UserRegistrationInput) error {
// 	// First, check if the email already exists
//     emailFilter := bson.M{"email": input.Email}
//     emailCount, err := us.collection.CountDocuments(ctx, emailFilter)
//     if err != nil {
//         return fmt.Errorf("error checking if email already exists: %w", err)
//     }
//     if emailCount > 0 {
//         return ErrUserAlreadyExists
//     }

//     // Then, check if the username already exists
//     usernameFilter := bson.M{"profileInformation.username": input.ProfileInformation.Username}
//     usernameCount, err := us.collection.CountDocuments(ctx, usernameFilter)
//     if err != nil {
//         return fmt.Errorf("error checking if username already exists: %w", err)
//     }
//     if usernameCount > 0 {
//         return ErrUsernameAlreadyTaken
//     }

// 	// Create a new user
// 	user, err := models.NewUserfromInput(input)
// 	if err != nil {
// 		return fmt.Errorf("error creating user from input: %w", err)
// 	}
	
// 	// Insert the user into the database
// 	_, err = us.collection.InsertOne(ctx, user)
// 	if err != nil {
// 		return fmt.Errorf("error inserting user into database: %w", err)
// 	}

// 	return nil
// }

// func (us *UserService) GetUserByEmail(ctx context.Context, email, password string) (*models.User, error) {
// 	var user models.User
// 	filter := bson.M{"email": email}

// 	result := us.collection.FindOne(ctx, filter)
// 	if result.Err() != nil {
// 		return &models.User{}, fmt.Errorf("error fetching user with email %s: %w", email, result.Err())
// 	}

// 	err := result.Decode(&user)
// 	if err != nil {
// 		return &models.User{}, fmt.Errorf("error decoding user data for email %s: %w", email, err)
// 	}

// 	if !us.hasher.CheckPasswordHash(password, user.PasswordHash) {
// 		return &models.User{}, ErrInvalidUserCredentials
// 	}	

// 	return &user, nil
// }