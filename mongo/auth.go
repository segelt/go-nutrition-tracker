package mongo

import (
	"context"
	"fmt"
	"nutritiontracker/internal"
	"nutritiontracker/resource"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ resource.AuthService = (*AuthService)(nil)

type AuthService struct {
	db *DB
}

func NewAuthService(db *DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) RegisterUser(create resource.UserInput) error {
	newUser := resource.User{
		ID:       primitive.NewObjectID(),
		Username: create.Username,
		Password: create.Password,
	}
	err := newUser.HashPassword(create.Password)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.db.operationTimeout)
	defer cancel()

	coll := s.db.client.Database("nutrition-tracker").Collection("users")
	_, err = coll.InsertOne(ctx, newUser)
	return err
}

func (s *AuthService) LoginUser(login resource.UserInput) (generatedToken string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), s.db.operationTimeout)
	defer cancel()

	coll := s.db.client.Database("nutrition-tracker").Collection("users")
	filter := bson.D{{"username", login.Username}}

	var targetUser resource.User
	err = coll.FindOne(ctx, filter).Decode(&targetUser)

	if err != nil {
		return "", err
	}

	err = targetUser.CheckPassword(login.Password)
	if err != nil {
		return "", fmt.Errorf("invalid password")
	}

	token, err := internal.GenerateToken(targetUser.ID.Hex(), targetUser.Username)
	if err != nil {
		return "", fmt.Errorf("failed to generate token for the user after succesful login", err)
	}

	return token, nil
}
