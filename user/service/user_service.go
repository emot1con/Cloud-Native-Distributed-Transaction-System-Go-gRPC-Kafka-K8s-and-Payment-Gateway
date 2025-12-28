package service

import (
	"fmt"
	"os"
	"strconv"
	"time"
	"user_service/auth"
	"user_service/proto"
	"user_service/repository"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

type UserService struct {
	validator *validator.Validate
	repo      repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		validator: validator.New(),
		repo:      repo,
	}
}

func (u *UserService) Register(payload *proto.RegisterPayload) error {
	logrus.Info("get user")
	if user, err := u.repo.GetUserByEmail(payload.Email); err == nil && user != nil {
		return fmt.Errorf("email already exists")
	}

	if payload.Provider == "" || payload.Provider == "local" {
		hashedPassword, err := auth.GeneratePasswordHash(payload.Password)
		if err != nil {
			return err
		}
		payload.Password = hashedPassword
	} else {
		payload.Password = ""
	}

	logrus.Info("registering user to db")
	if err := u.repo.Register(payload); err != nil {
		return err
	}

	return nil
}

func (u *UserService) Login(payload *proto.LoginPayload) (string, string, error) {
	user, err := u.repo.GetUserByEmail(payload.Email)
	if err != nil {
		return "", "", fmt.Errorf("email not found")
	}

	if err := auth.ComparePasswordHash(user.Password, []byte(payload.Password)); err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	token, err := auth.CreateJWTToken(int(user.ID), "user", 0, 0, 1)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := auth.CreateJWTToken(int(user.ID), "user", 0, 3, 0)
	if err != nil {
		return "", "", err
	}
	logrus.Info("logged!")

	return token, refreshToken, nil
}

func (u *UserService) RefreshToken(payload string) (string, string, error) {
	token, err := jwt.Parse(payload, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return "", "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", "", fmt.Errorf("invalid claims")
	}

	expiredAt := int64(claims["exp"].(float64))
	if time.Now().Unix() > expiredAt {
		return "", "", fmt.Errorf("refresh token expired")
	}

	ID, err := strconv.Atoi(claims["userID"].(string))
	if err != nil {
		return "", "", fmt.Errorf("failed to parse user ID")
	}

	newToken, err := auth.CreateJWTToken(ID, "user", 0, 0, 1)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err := auth.CreateJWTToken(ID, "user", 0, 3, 0)
	if err != nil {
		return "", "", err
	}

	return newToken, newRefreshToken, nil
}

func (u *UserService) GetUserByID(ID int) (*proto.User, error) {
	user, err := u.repo.GetUserByID(ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) GetUserByEmail(email string) (*proto.User, error) {
	user, err := u.repo.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserService) UpdateUser(user *proto.User) error {
	if err := u.repo.UpdateUser(user); err != nil {
		return err
	}
	return nil
}
