package service

import (
	"context"

	"github.com/Talan-Application/quiz-service/internal/domain"
	"github.com/Talan-Application/quiz-service/internal/repository"
	"go.uber.org/zap"
)

type QuestionService struct {
	repo   repository.QuestionRepository
	logger *zap.Logger
}

func NewQuestionService(repo repository.QuestionRepository, logger *zap.Logger) *QuestionService {
	return &QuestionService{repo: repo, logger: logger}
}

func (s *QuestionService) Create(ctx context.Context, question *domain.Question) (*domain.Question, error) {
	created, err := s.repo.Create(ctx, question)
	if err != nil {
		s.logger.Error("failed to create question", zap.Error(err))
		return nil, err
	}
	return created, nil
}

func (s *QuestionService) GetByID(ctx context.Context, id int64) (*domain.Question, error) {
	question, err := s.repo.GetById(ctx, id)
	if err != nil {
		s.logger.Error("failed to get question", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return question, nil
}

func (s *QuestionService) GetAll(ctx context.Context, quizID int64, limit *int, offset *int) ([]domain.Question, error) {
	questions, err := s.repo.GetAll(ctx, quizID, limit, offset)
	if err != nil {
		s.logger.Error("failed to get questions", zap.Int64("quiz_id", quizID), zap.Error(err))
		return nil, err
	}
	return questions, nil
}

func (s *QuestionService) Update(ctx context.Context, id int64, question *domain.Question) (*domain.Question, error) {
	updated, err := s.repo.Update(ctx, id, question)
	if err != nil {
		s.logger.Error("failed to update question", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return updated, nil
}

func (s *QuestionService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete question", zap.Int64("id", id), zap.Error(err))
		return err
	}
	return nil
}
