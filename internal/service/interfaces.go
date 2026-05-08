package service

import (
	"context"

	"github.com/Talan-Application/quiz-service/internal/domain"
)

type IQuizService interface {
	Create(ctx context.Context, quiz *domain.Quiz) (*domain.Quiz, error)
	GetByID(ctx context.Context, id int64) (*domain.Quiz, error)
	GetAll(ctx context.Context, limit *int, offset *int) ([]domain.Quiz, error)
	Update(ctx context.Context, id int64, quiz *domain.Quiz) (*domain.Quiz, error)
	Delete(ctx context.Context, id int64) error
}