package handler

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	quizv1 "github.com/Talan-Application/proto-generation/quiz/v1"
	"github.com/Talan-Application/quiz-service/internal/domain"
)

func (h *Handler) CreateQuestionWithAnswers(ctx context.Context, req *quizv1.CreateQuestionWithAnswersRequest) (*quizv1.QuestionWithAnswersResponse, error) {
	created, err := h.questionSvc.CreateWithAnswers(ctx, req)
	if err != nil {
		return nil, h.toQuestionGRPCError(err)
	}
	return questionWithAnswersToProto(created), nil
}

func (h *Handler) GetQuestion(ctx context.Context, req *quizv1.GetQuestionRequest) (*quizv1.QuestionWithAnswersResponse, error) {
	result, err := h.questionSvc.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, h.toQuestionGRPCError(err)
	}
	return questionWithAnswersToProto(result), nil
}

func (h *Handler) GetAllQuestions(ctx context.Context, req *quizv1.GetAllQuestionsRequest) (*quizv1.GetAllQuestionsWithAnswersResponse, error) {
	results, err := h.questionSvc.GetAll(ctx, req)
	if err != nil {
		return nil, h.toQuestionGRPCError(err)
	}

	protos := make([]*quizv1.QuestionWithAnswersResponse, len(results))
	for i := range results {
		protos[i] = questionWithAnswersToProto(&results[i])
	}

	return &quizv1.GetAllQuestionsWithAnswersResponse{Questions: protos}, nil
}

func (h *Handler) UpdateQuestionWithAnswers(ctx context.Context, req *quizv1.UpdateQuestionWithAnswersRequest) (*quizv1.QuestionWithAnswersResponse, error) {
	result, err := h.questionSvc.UpdateWithAnswers(ctx, req.GetId(), req)
	if err != nil {
		return nil, h.toQuestionGRPCError(err)
	}
	return questionWithAnswersToProto(result), nil
}

func (h *Handler) DeleteQuestion(ctx context.Context, req *quizv1.DeleteQuestionRequest) (*quizv1.DeleteQuestionResponse, error) {
	if err := h.questionSvc.Delete(ctx, req.GetId()); err != nil {
		return nil, h.toQuestionGRPCError(err)
	}
	return &quizv1.DeleteQuestionResponse{Message: "question deleted"}, nil
}

func questionWithAnswersToProto(q *domain.QuestionWithAnswers) *quizv1.QuestionWithAnswersResponse {
	answers := make([]*quizv1.AnswerResponse, len(q.Answers))
	for i := range q.Answers {
		answers[i] = answerToProto(&q.Answers[i])
	}
	return &quizv1.QuestionWithAnswersResponse{
		Id:             q.ID,
		QuizId:         q.QuizID,
		Text:           q.Text,
		Context:        q.Context,
		VideoAnswerUrl: q.VideoAnswerUrl,
		Order:          q.Order,
		CreatedAt:      q.CreatedAt.Unix(),
		UpdatedAt:      q.UpdatedAt.Unix(),
		Answers:        answers,
	}
}

func (h *Handler) toQuestionGRPCError(err error) error {
	if errors.Is(err, domain.ErrQuestionNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	h.log.Error("unexpected error", zap.Error(err))
	return status.Error(codes.Internal, "internal error")
}
