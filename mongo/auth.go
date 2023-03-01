package mongo

import (
	"context"
	"html"
	"nutritiontracker/resource"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
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
	err := processUserInfoBeforeRegistration(&newUser)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.db.operationTimeout)
	defer cancel()

	coll := s.db.client.Database("nutrition-tracker").Collection("users")
	_, err = coll.InsertOne(ctx, newUser)
	return err
}

func processUserInfoBeforeRegistration(u *resource.User) error {
	//turn password into hash
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	//remove spaces in username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))

	return nil
}
