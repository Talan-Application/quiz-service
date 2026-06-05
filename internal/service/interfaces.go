package service

import (
	"context"

	quizv1 "github.com/Talan-Application/proto-generation/quiz/v1"
	"github.com/Talan-Application/quiz-service/internal/domain"
)

type IQuizService interface {
	Create(ctx context.Context, req *quizv1.CreateQuizRequest, userID int64) (*domain.Quiz, error)
	GetByID(ctx context.Context, id int64) (*domain.Quiz, error)
	GetAll(ctx context.Context, req *quizv1.GetAllQuizzesRequest) ([]domain.Quiz, error)
	GetAllByAuthor(ctx context.Context, req *quizv1.GetMyQuizzesRequest, userID int64) ([]domain.Quiz, error)
	Publish(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, req *quizv1.UpdateQuizRequest) (*domain.Quiz, error)
	Delete(ctx context.Context, id int64) error
}

type IQuestionService interface {
	CreateWithAnswers(ctx context.Context, req *quizv1.CreateQuestionWithAnswersRequest) (*domain.QuestionWithAnswers, error)
	GetByID(ctx context.Context, id int64) (*domain.QuestionWithAnswers, error)
	GetAll(ctx context.Context, req *quizv1.GetAllQuestionsRequest) ([]domain.QuestionWithAnswers, error)
	UpdateWithAnswers(ctx context.Context, id int64, req *quizv1.UpdateQuestionWithAnswersRequest) (*domain.QuestionWithAnswers, error)
	Delete(ctx context.Context, id int64) error
}

type IQuizResultService interface {
	Submit(ctx context.Context, userID int64, req *quizv1.SubmitQuizRequest) (*domain.QuizResult, error)
	GetResults(ctx context.Context, quizID, userID int64) ([]domain.QuizResult, error)
}
