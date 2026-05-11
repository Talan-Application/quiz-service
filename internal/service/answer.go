package service

import (
	"context"

	"github.com/Talan-Application/quiz-service/internal/domain"
	"github.com/Talan-Application/quiz-service/internal/repository"
	"go.uber.org/zap"
)

type AnswerService struct {
	repo   repository.AnswerRepository
	logger *zap.Logger
}

func NewAnswerService(repo repository.AnswerRepository, logger *zap.Logger) *AnswerService {
	return &AnswerService{repo: repo, logger: logger}
}

func (s *AnswerService) Create(ctx context.Context, answer *domain.Answer) (*domain.Answer, error) {
	created, err := s.repo.Create(ctx, answer)
	if err != nil {
		s.logger.Error("failed to create answer", zap.Error(err))
		return nil, err
	}
	return created, nil
}

func (s *AnswerService) GetByID(ctx context.Context, id int64) (*domain.Answer, error) {
	answer, err := s.repo.GetById(ctx, id)
	if err != nil {
		s.logger.Error("failed to get answer", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return answer, nil
}

func (s *AnswerService) GetAll(ctx context.Context, questionID int64, limit *int, offset *int) ([]domain.Answer, error) {
	answers, err := s.repo.GetAll(ctx, questionID, limit, offset)
	if err != nil {
		s.logger.Error("failed to get answers", zap.Int64("question_id", questionID), zap.Error(err))
		return nil, err
	}
	return answers, nil
}

func (s *AnswerService) Update(ctx context.Context, id int64, answer *domain.Answer) (*domain.Answer, error) {
	updated, err := s.repo.Update(ctx, id, answer)
	if err != nil {
		s.logger.Error("failed to update answer", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return updated, nil
}

func (s *AnswerService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete answer", zap.Int64("id", id), zap.Error(err))
		return err
	}
	return nil
}
