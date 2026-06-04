package service

import (
	"context"

	"go.uber.org/zap"

	quizv1 "github.com/Talan-Application/proto-generation/quiz/v1"
	"github.com/Talan-Application/quiz-service/internal/domain"
	"github.com/Talan-Application/quiz-service/internal/repository"
)

type AnswerService struct {
	repo   repository.AnswerRepository
	logger *zap.Logger
}

func NewAnswerService(repo repository.AnswerRepository, logger *zap.Logger) *AnswerService {
	return &AnswerService{repo: repo, logger: logger}
}

func (s *AnswerService) Create(ctx context.Context, req *quizv1.CreateAnswerRequest) (*domain.Answer, error) {
	answer := &domain.Answer{
		QuestionID: req.GetQuestionId(),
		Text:       req.GetText(),
		Correct:    req.GetCorrect(),
	}
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

func (s *AnswerService) GetAll(ctx context.Context, req *quizv1.GetAllAnswersRequest) ([]domain.Answer, error) {
	var limit, offset *int
	if req.Limit != nil {
		v := int(req.GetLimit())
		limit = &v
	}
	if req.Offset != nil {
		v := int(req.GetOffset())
		offset = &v
	}

	answers, err := s.repo.GetAll(ctx, req.GetQuestionId(), limit, offset)
	if err != nil {
		s.logger.Error("failed to get answers", zap.Int64("question_id", req.GetQuestionId()), zap.Error(err))
		return nil, err
	}
	return answers, nil
}

func (s *AnswerService) Update(ctx context.Context, id int64, req *quizv1.UpdateAnswerRequest) (*domain.Answer, error) {
	answer := &domain.Answer{
		Text:    req.GetText(),
		Correct: req.GetCorrect(),
	}
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
