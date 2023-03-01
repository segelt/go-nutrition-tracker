package mongo

import (
	"fmt"
	"nutritiontracker/resource"
)

var _ resource.AuthService = (*AuthService)(nil)

type AuthService struct {
	db *DB
}

func NewAuthService(db *DB) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) RegisterUser(create resource.UserInput) error {
	return fmt.Errorf("not implemented")
}
