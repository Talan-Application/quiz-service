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

type IQuestionService interface {
	Create(ctx context.Context, question *domain.Question) (*domain.Question, error)
	GetByID(ctx context.Context, id int64) (*domain.Question, error)
	GetAll(ctx context.Context, quizID int64, limit *int, offset *int) ([]domain.Question, error)
	Update(ctx context.Context, id int64, question *domain.Question) (*domain.Question, error)
	Delete(ctx context.Context, id int64) error
}

type IAnswerService interface {
	Create(ctx context.Context, answer *domain.Answer) (*domain.Answer, error)
	GetByID(ctx context.Context, id int64) (*domain.Answer, error)
	GetAll(ctx context.Context, questionID int64, limit *int, offset *int) ([]domain.Answer, error)
	Update(ctx context.Context, id int64, answer *domain.Answer) (*domain.Answer, error)
	Delete(ctx context.Context, id int64) error
}

type IQuizResultService interface {
	Submit(ctx context.Context, result *domain.QuizResult) (*domain.QuizResult, error)
	GetByQuizAndUser(ctx context.Context, quizID, userID int64) ([]domain.QuizResult, error)
	GetByQuiz(ctx context.Context, quizID int64) ([]domain.QuizResult, error)
}