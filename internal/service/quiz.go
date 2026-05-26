package service

import (
	"context"

	"github.com/Talan-Application/quiz-service/internal/domain"
	"github.com/Talan-Application/quiz-service/internal/repository"
	"go.uber.org/zap"
)

type QuizService struct {
	repo   repository.QuizRepository
	logger *zap.Logger
}

func NewQuizService(repo repository.QuizRepository, logger *zap.Logger) *QuizService {
	return &QuizService{repo: repo, logger: logger}
}

func (s *QuizService) Create(ctx context.Context, quiz *domain.Quiz) (*domain.Quiz, error) {
	created, err := s.repo.Create(ctx, quiz)
	if err != nil {
		s.logger.Error("failed to create quiz", zap.Error(err))
		return nil, err
	}
	return created, nil
}

func (s *QuizService) GetByID(ctx context.Context, id int64) (*domain.Quiz, error) {
	quiz, err := s.repo.GetById(ctx, id)
	if err != nil {
		s.logger.Error("failed to get quiz", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return quiz, nil
}

func (s *QuizService) GetAll(ctx context.Context, status *domain.QuizStatus, limit *int, offset *int) ([]domain.Quiz, error) {
	quizzes, err := s.repo.GetAll(ctx, status, limit, offset)
	if err != nil {
		s.logger.Error("failed to get quizzes", zap.Error(err))
		return nil, err
	}
	return quizzes, nil
}

func (s *QuizService) Publish(ctx context.Context, id int64) (*domain.Quiz, error) {
	quiz, err := s.repo.Publish(ctx, id)
	if err != nil {
		s.logger.Error("failed to publish quiz", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return quiz, nil
}

func (s *QuizService) GetAllByAuthor(ctx context.Context, authorID int64, limit *int, offset *int) ([]domain.Quiz, error) {
	quizzes, err := s.repo.GetAllByAuthor(ctx, authorID, limit, offset)
	if err != nil {
		s.logger.Error("failed to get quizzes by author", zap.Int64("author_id", authorID), zap.Error(err))
		return nil, err
	}
	return quizzes, nil
}

func (s *QuizService) Update(ctx context.Context, id int64, quiz *domain.Quiz) (*domain.Quiz, error) {
	updated, err := s.repo.Update(ctx, id, quiz)
	if err != nil {
		s.logger.Error("failed to update quiz", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return updated, nil
}

func (s *QuizService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete quiz", zap.Int64("id", id), zap.Error(err))
		return err
	}
	return nil
}
