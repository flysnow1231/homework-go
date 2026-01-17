package repo

import (
	"blog/internal/model"
	"context"
	"gorm.io/gorm"
)

type PostRepo struct {
	db *gorm.DB
}

func NewPostRepo(db *gorm.DB) *PostRepo { return &PostRepo{db: db} }

func (r *PostRepo) Create(ctx context.Context, p *model.Post) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *PostRepo) GetByUsername(ctx context.Context, userId int, page int, size int) ([]model.Post, error) {
	var post []model.Post
	var total int64
	err := r.db.WithContext(ctx).
		Model(post).
		Where("user_id = ?", userId).
		Count(&total).Error

	if err != nil {
		return nil, err
	}

	offset := page*size - 1
	var list []model.Post
	err = r.db.WithContext(ctx).Model(post).Where("user_id = ?", userId).
		Order("created_at DESC").
		Limit(size).
		Offset(offset).
		Find(&list).
		Error
	return list, nil
}
