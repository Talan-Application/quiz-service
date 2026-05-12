package service

import (
	"context"

	"go.uber.org/zap"

	"github.com/Talan-Application/quiz-service/internal/domain"
	"github.com/Talan-Application/quiz-service/internal/repository"
)

type QuizResultService struct {
	repo   repository.QuizResultRepository
	logger *zap.Logger
}

func NewQuizResultService(repo repository.QuizResultRepository, logger *zap.Logger) *QuizResultService {
	return &QuizResultService{repo: repo, logger: logger}
}

func (s *QuizResultService) Submit(ctx context.Context, result *domain.QuizResult) (*domain.QuizResult, error) {
	saved, err := s.repo.Save(ctx, result)
	if err != nil {
		s.logger.Error("failed to save quiz result",
			zap.Int64("quiz_id", result.QuizID),
			zap.Int64("user_id", result.UserID),
			zap.Error(err),
		)
		return nil, err
	}
	return saved, nil
}

func (s *QuizResultService) GetByQuizAndUser(ctx context.Context, quizID, userID int64) ([]domain.QuizResult, error) {
	results, err := s.repo.GetByQuizAndUser(ctx, quizID, userID)
	if err != nil {
		s.logger.Error("failed to get quiz results",
			zap.Int64("quiz_id", quizID),
			zap.Int64("user_id", userID),
			zap.Error(err),
		)
		return nil, err
	}
	return results, nil
}

func (s *QuizResultService) GetByQuiz(ctx context.Context, quizID int64) ([]domain.QuizResult, error) {
	results, err := s.repo.GetByQuiz(ctx, quizID)
	if err != nil {
		s.logger.Error("failed to get quiz results", zap.Int64("quiz_id", quizID), zap.Error(err))
		return nil, err
	}
	return results, nil
}
