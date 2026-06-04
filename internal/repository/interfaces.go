package repository

import (
	"context"

	"github.com/Talan-Application/quiz-service/internal/domain"
)

type QuizRepository interface {
	Create(ctx context.Context, quiz *domain.Quiz) (*domain.Quiz, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, quiz *domain.Quiz) (*domain.Quiz, error)
	GetAll(ctx context.Context, status *domain.QuizStatus, limit *int, offset *int) ([]domain.Quiz, error)
	GetById(ctx context.Context, id int64) (*domain.Quiz, error)
	Publish(ctx context.Context, id int64) error
	GetAllByAuthor(ctx context.Context, authorID int64, limit *int, offset *int) ([]domain.Quiz, error)
}

type QuestionRepository interface {
	Create(ctx context.Context, question *domain.Question) (*domain.Question, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, question *domain.Question) (*domain.Question, error)
	GetAll(ctx context.Context, quizID int64, limit *int, offset *int) ([]domain.Question, error)
	GetById(ctx context.Context, id int64) (*domain.Question, error)
}

type AnswerRepository interface {
	Create(ctx context.Context, answer *domain.Answer) (*domain.Answer, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, answer *domain.Answer) (*domain.Answer, error)
	GetAll(ctx context.Context, questionID int64, limit *int, offset *int) ([]domain.Answer, error)
	GetById(ctx context.Context, id int64) (*domain.Answer, error)
}

type QuizResultRepository interface {
	Save(ctx context.Context, result *domain.QuizResult) (*domain.QuizResult, error)
	GetByQuizAndUser(ctx context.Context, quizID, userID int64) ([]domain.QuizResult, error)
	GetByQuiz(ctx context.Context, quizID int64) ([]domain.QuizResult, error)
}
