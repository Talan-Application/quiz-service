package repository

import (
	"context"

	"github.com/Talan-Application/quiz-service/internal/domain"
)

type QuizRepository interface {
	Create(ctx context.Context, quiz *domain.Quiz) (*domain.Quiz, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, quiz *domain.Quiz) (*domain.Quiz, error)
	GetAll(ctx context.Context, limit *int, offset *int) ([]domain.Quiz, error)
	GetById(ctx context.Context, id int64) (*domain.Quiz, error)
}
