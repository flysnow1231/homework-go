package service

import (
	"blog/internal/model"
	"blog/internal/mq"
	"blog/internal/repo"
	"context"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repo.UserRepo
	mq   *mq.Producer
	log  *zap.Logger
}

func NewUserService(repo *repo.UserRepo, producer *mq.Producer, log *zap.Logger) *UserService {
	return &UserService{repo: repo, mq: producer, log: log}
}

func (s *UserService) CreateUser(ctx context.Context, u *model.User) (*model.User, error) {
	s.log.Info("create user", zap.String("create user", u.Username), zap.String("Email", u.Email))
	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserService) GetUser(ctx context.Context, username string, pwd string) (*model.User, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil || user == nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pwd))
	if err != nil {
		return nil, fmt.Errorf("query user %v failed: %w", username, err)
	}
	return user, nil

}
