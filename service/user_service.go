package service

import (
	"context"
	"fmt"
	"os"
	"time"
	"transaction_api/model/user"
	"transaction_api/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	Repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		Repo: repo,
	}
}

func (s *UserService) GetUserDetail(ctx context.Context, username string) (*user.User, error) {
	uf, err := s.Repo.GetInquiryUserRepository(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return uf, nil
}
func (s *UserService) GetUserService(ctx context.Context, username, password string) (*user.User, error) {
	uf, err := s.Repo.GetInquiryUserRepository(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	if !checkPassword(password, uf.Password) {
		return nil, fmt.Errorf("password did not match")
	}

	token, err := generateToken(uf.Id, uf.Email, uf.Role)
	if err != nil {
		return nil, fmt.Errorf("generate user token: %w", err)
	}

	_, err = s.Repo.UpdateUserRepository(ctx, uf.Id, username, token)
	if err != nil {
		return nil, fmt.Errorf("update user token: %w", err)
	}

	return &user.User{
		Name:  uf.Name,
		Email: uf.Email,
		Token: &token,
	}, nil
}

func checkPassword(inputPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword(
		[]byte(hashedPassword), // hash yang tersimpan di database
		[]byte(inputPassword),  // password yang diinput user saat login
	)
	return err == nil // nil artinya cocok, error artinya tidak cocok
}

func generateToken(userID int, email, role string) (string, error) {
	// payload / claims -> mirip object pertama di jwt.sign()
	claims := jwt.MapClaims{
		"id":    userID,
		"email": email,
		"role":  role,
		"exp":   time.Now().Add(1 * time.Hour).Unix(), // mirip expiresIn: "1d"
	}

	// buat token dengan algoritma HS256 (default jsonwebtoken juga HS256)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign pakai secret key -> mirip parameter ke-2 di jwt.sign()
	secret := os.Getenv("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
