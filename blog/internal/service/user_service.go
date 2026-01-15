package service

import (
	"blog/internal/model"
	"blog/internal/mq"
	"blog/internal/repo"
	"context"
)

type UserService struct {
	repo *repo.UserRepo
	mq   *mq.Producer
}

func NewUserService(repo *repo.UserRepo, producer *mq.Producer) *UserService {
	return &UserService{repo: repo, mq: producer}
}

func (s *UserService) CreateUser(ctx context.Context, u *model.User) (*model.User, error) {

	if err := s.repo.Create(ctx, u); err != nil {
		return nil, err
	}

	return u, nil
}

func (s *UserService) GetUser(ctx context.Context, username string) (*model.User, error) {
	return s.repo.GetByUsername(ctx, username)
}
