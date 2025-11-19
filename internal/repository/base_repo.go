package repository

import (
	"context"

	"gorm.io/gorm"
)

// Struct implementasi
type baseRepository[T any, K comparable] struct {
	db *gorm.DB
}

// Constructor helper
func NewBaseRepository[T any, K comparable](db *gorm.DB) BaseRepository[T, K] {
	return &baseRepository[T, K]{db: db}
}

func (r *baseRepository[T, K]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *baseRepository[T, K]) GetByID(ctx context.Context, id K) (*T, error) {
	var entity T
	// GORM cukup pintar untuk tahu kolom primary key,
	// asalkan struct T punya tag `gorm:"primaryKey"`
	if err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *baseRepository[T, K]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *baseRepository[T, K]) Delete(ctx context.Context, id K) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, "id = ?", id).Error
}