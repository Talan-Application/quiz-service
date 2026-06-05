package service

import (
	"context"

	"go.uber.org/zap"

	quizv1 "github.com/Talan-Application/proto-generation/quiz/v1"
	"github.com/Talan-Application/quiz-service/internal/domain"
	"github.com/Talan-Application/quiz-service/internal/repository"
)

type QuestionService struct {
	repo       repository.QuestionRepository
	answerRepo repository.AnswerRepository
	logger     *zap.Logger
}

func NewQuestionService(repo repository.QuestionRepository, answerRepo repository.AnswerRepository, logger *zap.Logger) *QuestionService {
	return &QuestionService{repo: repo, answerRepo: answerRepo, logger: logger}
}

func (s *QuestionService) CreateWithAnswers(ctx context.Context, req *quizv1.CreateQuestionWithAnswersRequest) (*domain.QuestionWithAnswers, error) {
	question := &domain.Question{
		QuizID:         req.GetQuizId(),
		Text:           req.GetText(),
		Context:        req.GetContext(),
		VideoAnswerUrl: req.GetVideoAnswerUrl(),
		Order:          req.GetOrder(),
	}

	answers := make([]domain.Answer, len(req.GetAnswers()))
	for i, a := range req.GetAnswers() {
		answers[i] = domain.Answer{
			Text:    a.GetText(),
			Correct: a.GetCorrect(),
		}
	}

	result, err := s.repo.CreateWithAnswers(ctx, question, answers)
	if err != nil {
		s.logger.Error("failed to create question with answers", zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (s *QuestionService) GetByID(ctx context.Context, id int64) (*domain.QuestionWithAnswers, error) {
	result, err := s.repo.GetByIdWithAnswers(ctx, id)
	if err != nil {
		s.logger.Error("failed to get question with answers", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (s *QuestionService) GetAll(ctx context.Context, req *quizv1.GetAllQuestionsRequest) ([]domain.QuestionWithAnswers, error) {
	var limit, offset *int
	if req.Limit != nil {
		v := int(req.GetLimit())
		limit = &v
	}
	if req.Offset != nil {
		v := int(req.GetOffset())
		offset = &v
	}

	results, err := s.repo.GetAllWithAnswers(ctx, req.GetQuizId(), limit, offset)
	if err != nil {
		s.logger.Error("failed to get questions with answers", zap.Int64("quiz_id", req.GetQuizId()), zap.Error(err))
		return nil, err
	}
	return results, nil
}

func (s *QuestionService) UpdateWithAnswers(ctx context.Context, id int64, req *quizv1.UpdateQuestionWithAnswersRequest) (*domain.QuestionWithAnswers, error) {
	question := &domain.Question{
		Text:           req.GetText(),
		Context:        req.GetContext(),
		VideoAnswerUrl: req.GetVideoAnswerUrl(),
		Order:          req.GetOrder(),
	}

	answers := make([]domain.Answer, len(req.GetAnswers()))
	for i, a := range req.GetAnswers() {
		answers[i] = domain.Answer{
			Text:    a.GetText(),
			Correct: a.GetCorrect(),
		}
	}

	result, err := s.repo.UpdateWithAnswers(ctx, id, question, answers)
	if err != nil {
		s.logger.Error("failed to update question with answers", zap.Int64("id", id), zap.Error(err))
		return nil, err
	}
	return result, nil
}

func (s *QuestionService) Delete(ctx context.Context, id int64) error {
	if err := s.answerRepo.DeleteByQuestionID(ctx, id); err != nil {
		s.logger.Error("failed to delete answers for question", zap.Int64("id", id), zap.Error(err))
		return err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		s.logger.Error("failed to delete question", zap.Int64("id", id), zap.Error(err))
		return err
	}
	return nil
}
