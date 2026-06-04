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

func (h *Handler) CreateQuestion(ctx context.Context, req *quizv1.CreateQuestionRequest) (*quizv1.QuestionResponse, error) {
	created, err := h.questionSvc.Create(ctx, req)
	if err != nil {
		return nil, h.toQuestionGRPCError(err)
	}

	return questionToProto(created), nil
}

func (h *Handler) GetQuestion(ctx context.Context, req *quizv1.GetQuestionRequest) (*quizv1.QuestionResponse, error) {
	question, err := h.questionSvc.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, h.toQuestionGRPCError(err)
	}

	return questionToProto(question), nil
}

func (h *Handler) GetAllQuestions(ctx context.Context, req *quizv1.GetAllQuestionsRequest) (*quizv1.GetAllQuestionsResponse, error) {
	questions, err := h.questionSvc.GetAll(ctx, req)
	if err != nil {
		return nil, h.toQuestionGRPCError(err)
	}

	protos := make([]*quizv1.QuestionResponse, len(questions))
	for i := range questions {
		protos[i] = questionToProto(&questions[i])
	}

	return &quizv1.GetAllQuestionsResponse{Questions: protos}, nil
}

func (h *Handler) UpdateQuestion(ctx context.Context, req *quizv1.UpdateQuestionRequest) (*quizv1.QuestionResponse, error) {
	updated, err := h.questionSvc.Update(ctx, req.GetId(), req)
	if err != nil {
		return nil, h.toQuestionGRPCError(err)
	}

	return questionToProto(updated), nil
}

func (h *Handler) DeleteQuestion(ctx context.Context, req *quizv1.DeleteQuestionRequest) (*quizv1.DeleteQuestionResponse, error) {
	if err := h.questionSvc.Delete(ctx, req.GetId()); err != nil {
		return nil, h.toQuestionGRPCError(err)
	}

	return &quizv1.DeleteQuestionResponse{Message: "question deleted"}, nil
}

func questionToProto(q *domain.Question) *quizv1.QuestionResponse {
	return &quizv1.QuestionResponse{
		Id:             q.ID,
		QuizId:         q.QuizID,
		Text:           q.Text,
		Context:        q.Context,
		VideoAnswerUrl: q.VideoAnswerUrl,
		Order:          q.Order,
		CreatedAt:      q.CreatedAt.Unix(),
		UpdatedAt:      q.UpdatedAt.Unix(),
	}
}

func (h *Handler) toQuestionGRPCError(err error) error {
	if errors.Is(err, domain.ErrQuestionNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	h.log.Error("unexpected error", zap.Error(err))
	return status.Error(codes.Internal, "internal error")
}
