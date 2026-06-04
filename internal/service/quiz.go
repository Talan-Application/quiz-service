package service

import (
	"context"

	"go.uber.org/zap"

	quizv1 "github.com/Talan-Application/proto-generation/quiz/v1"
	"github.com/Talan-Application/quiz-service/internal/domain"
	"github.com/Talan-Application/quiz-service/internal/repository"
)

type QuizService struct {
	repo   repository.QuizRepository
	logger *zap.Logger
}

func NewQuizService(repo repository.QuizRepository, logger *zap.Logger) *QuizService {
	return &QuizService{repo: repo, logger: logger}
}

func (s *QuizService) Create(ctx context.Context, req *quizv1.CreateQuizRequest, userID int64) (*domain.Quiz, error) {
	quiz := &domain.Quiz{
		Title:           req.GetTitle(),
		Language:        req.GetLanguage(),
		AuthorID:        userID,
		Type:            domain.QuizType(req.GetType()),
		CommonSubjectID: req.GetCommonSubjectId(),
		IsEntStandard:   req.GetIsEntStandard(),
		Status:          domain.QuizStatusDraft,
	}
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

func (s *QuizService) GetAll(ctx context.Context, req *quizv1.GetAllQuizzesRequest) ([]domain.Quiz, error) {
	var limit, offset *int
	var status *domain.QuizStatus

	if req.Limit != nil {
		v := int(req.GetLimit())
		limit = &v
	}
	if req.Offset != nil {
		v := int(req.GetOffset())
		offset = &v
	}
	if req.Status != nil {
		s := domain.QuizStatus(req.GetStatus())
		status = &s
	}

	quizzes, err := s.repo.GetAll(ctx, status, limit, offset)
	if err != nil {
		s.logger.Error("failed to get quizzes", zap.Error(err))
		return nil, err
	}
	return quizzes, nil
}

func (s *QuizService) Publish(ctx context.Context, id int64) error {
	if err := s.repo.Publish(ctx, id); err != nil {
		s.logger.Error("failed to publish quiz", zap.Int64("id", id), zap.Error(err))
		return err
	}
	return nil
}

func (s *QuizService) GetAllByAuthor(ctx context.Context, req *quizv1.GetMyQuizzesRequest, userID int64) ([]domain.Quiz, error) {
	var limit, offset *int
	if req.Limit != nil {
		v := int(req.GetLimit())
		limit = &v
	}
	if req.Offset != nil {
		v := int(req.GetOffset())
		offset = &v
	}

	quizzes, err := s.repo.GetAllByAuthor(ctx, userID, limit, offset)
	if err != nil {
		s.logger.Error("failed to get quizzes by author", zap.Int64("author_id", userID), zap.Error(err))
		return nil, err
	}
	return quizzes, nil
}

func (s *QuizService) Update(ctx context.Context, id int64, req *quizv1.UpdateQuizRequest) (*domain.Quiz, error) {
	quiz := &domain.Quiz{
		Title:         req.GetTitle(),
		Language:      req.GetLanguage(),
		Type:          domain.QuizType(req.GetType()),
		IsEntStandard: req.GetIsEntStandard(),
	}
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
