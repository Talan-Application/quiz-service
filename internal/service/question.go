package service

import (
	"context"

	"go.uber.org/zap"

	quizv1 "github.com/Talan-Application/proto-generation/quiz/v1"
	"github.com/Talan-Application/quiz-service/internal/domain"
	"github.com/Talan-Application/quiz-service/internal/repository"
)

type QuestionService struct {
	repo   repository.QuestionRepository
	logger *zap.Logger
}

func NewQuestionService(repo repository.QuestionRepository, logger *zap.Logger) *QuestionService {
	return &QuestionService{repo: repo, logger: logger}
}

func (s *QuestionService) Create(ctx context.Context, req *quizv1.CreateQuestionRequest) (*domain.Question, error) {
	question := &domain.Question{
		QuizID:         req.GetQuizId(),
		Text:           req.GetText(),
		Context:        req.GetContext(),
		VideoAnswerUrl: req.GetVideoAnswerUrl(),
		Order:          req.GetOrder(),
	}
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

func (s *QuestionService) GetAll(ctx context.Context, req *quizv1.GetAllQuestionsRequest) ([]domain.Question, error) {
	var limit, offset *int
	if req.Limit != nil {
		v := int(req.GetLimit())
		limit = &v
	}
	if req.Offset != nil {
		v := int(req.GetOffset())
		offset = &v
	}

	questions, err := s.repo.GetAll(ctx, req.GetQuizId(), limit, offset)
	if err != nil {
		s.logger.Error("failed to get questions", zap.Int64("quiz_id", req.GetQuizId()), zap.Error(err))
		return nil, err
	}
	return questions, nil
}

func (s *QuestionService) Update(ctx context.Context, id int64, req *quizv1.UpdateQuestionRequest) (*domain.Question, error) {
	question := &domain.Question{
		Text:           req.GetText(),
		Context:        req.GetContext(),
		VideoAnswerUrl: req.GetVideoAnswerUrl(),
		Order:          req.GetOrder(),
	}
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
