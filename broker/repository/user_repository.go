package repository

import (
	"broker/proto"
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type UserRepository interface {
	Register(*proto.RegisterRequest) error
	Login(*proto.LoginRequest) (*proto.TokenResponse, error)
	RefreshToken(*proto.RefreshTokenRequest) (*proto.TokenResponse, error)
	GetUserByID(ID int) (*proto.User, error)
	GoogleOauth() (*proto.URLResponse, error)
	FacebookOauth() (*proto.URLResponse, error)
	GithubOauth() (*proto.URLResponse, error)
	GoogleOauthCallback(code string) (*proto.TokenResponse, error)
	FacebookOauthCallback(code string) (*proto.TokenResponse, error)
	GithubOauthCallback(code string) (*proto.TokenResponse, error)
}

type UserRepositoryImpl struct {
	client proto.AuthServiceClient
}

func NewUserRepository() *UserRepositoryImpl {
	addr := os.Getenv("USER_SERVICE_ADDR")
	if addr == "" {
		logrus.Info("USER_SERVICE_ADDR not set, using default address")
		addr = "user-service:50001"
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Failed to connect to user service: %v", err)
	}
	logrus.Infof("Connected to user service at %s", addr)

	client := proto.NewAuthServiceClient(conn)

	return &UserRepositoryImpl{client: client}
}

func (u *UserRepositoryImpl) Register(payload *proto.RegisterRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := u.client.Register(ctx, &proto.RegisterRequest{
		Payload: &proto.RegisterPayload{
			FullName: payload.Payload.FullName,
			Email:    payload.Payload.Email,
			Password: payload.Payload.Password,
		},
	})
	if err != nil {
		logrus.Error("Failed to register: ", err)
		return err
	}

	logrus.Info("Registered successfully")
	return nil
}

func (u *UserRepositoryImpl) Login(payload *proto.LoginRequest) (*proto.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.client.Login(ctx, &proto.LoginRequest{
		Payload: &proto.LoginPayload{
			Email:    payload.Payload.Email,
			Password: payload.Payload.Password,
		},
	})
}

func (u *UserRepositoryImpl) RefreshToken(payload *proto.RefreshTokenRequest) (*proto.TokenResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.client.RefreshToken(ctx, &proto.RefreshTokenRequest{
		Payload: &proto.RefreshTokenPayload{
			RefreshToken: payload.Payload.RefreshToken,
		},
	})
}

func (u *UserRepositoryImpl) GetUserByID(ID int) (*proto.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return u.client.GetUserByID(ctx, &proto.GetUserRequest{
		UserId: int64(ID)},
	)
}

func (u *UserRepositoryImpl) GoogleOauth() (*proto.URLResponse, error) {
	return u.client.GoogleOauth(
		context.Background(), &proto.EmptyRequest{},
	)
}

func (u *UserRepositoryImpl) FacebookOauth() (*proto.URLResponse, error) {
	return u.client.FacebookOauth(
		context.Background(), &proto.EmptyRequest{},
	)
}

func (u *UserRepositoryImpl) GithubOauth() (*proto.URLResponse, error) {
	return u.client.GithubOauth(
		context.Background(), &proto.EmptyRequest{},
	)
}

func (u *UserRepositoryImpl) GoogleOauthCallback(code string) (*proto.TokenResponse, error) {
	return u.client.GoogleOauthCallback(
		context.Background(), &proto.CodeRequest{
			Code: code,
		},
	)
}

func (u *UserRepositoryImpl) FacebookOauthCallback(code string) (*proto.TokenResponse, error) {
	return u.client.FacebookOauthCallback(
		context.Background(), &proto.CodeRequest{
			Code: code,
		},
	)
}

func (u *UserRepositoryImpl) GithubOauthCallback(code string) (*proto.TokenResponse, error) {
	return u.client.GithubOauthCallback(
		context.Background(), &proto.CodeRequest{
			Code: code,
		},
	)
}
