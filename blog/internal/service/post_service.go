package service

import (
	"blog/internal/model"
	"blog/internal/mq"
	"blog/internal/repo"
	"context"
	"go.uber.org/zap"
)

type PostService struct {
	repo *repo.PostRepo
	mq   *mq.Producer
	log  *zap.Logger
}

func NewPostService(repo *repo.PostRepo, producer *mq.Producer, log *zap.Logger) *PostService {
	return &PostService{repo: repo, mq: producer, log: log}
}

func (s *PostService) CreatePost(ctx context.Context, post *model.Post) (*model.Post, error) {
	s.log.Info("create post", zap.Any("create post", post.ID))
	if err := s.repo.Create(ctx, post); err != nil {
		s.log.Error("create post", zap.Any("create post", post), zap.Error(err))
		return nil, err
	}
	return post, nil
}

func (s *PostService) QueryPosts(ctx context.Context, userId int, page int, size int) ([]model.Post, error) {
	posts, err := s.repo.GetByUsername(ctx, userId, page, size)
	if err != nil || posts == nil {
		return nil, err
	}

	return posts, nil
}
