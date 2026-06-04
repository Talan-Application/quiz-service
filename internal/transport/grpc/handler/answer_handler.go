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

func (h *Handler) CreateAnswer(ctx context.Context, req *quizv1.CreateAnswerRequest) (*quizv1.AnswerResponse, error) {
	created, err := h.answerSvc.Create(ctx, req)
	if err != nil {
		return nil, h.toAnswerGRPCError(err)
	}

	return answerToProto(created), nil
}

func (h *Handler) GetAnswer(ctx context.Context, req *quizv1.GetAnswerRequest) (*quizv1.AnswerResponse, error) {
	answer, err := h.answerSvc.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, h.toAnswerGRPCError(err)
	}

	return answerToProto(answer), nil
}

func (h *Handler) GetAllAnswers(ctx context.Context, req *quizv1.GetAllAnswersRequest) (*quizv1.GetAllAnswersResponse, error) {
	answers, err := h.answerSvc.GetAll(ctx, req)
	if err != nil {
		return nil, h.toAnswerGRPCError(err)
	}

	protos := make([]*quizv1.AnswerResponse, len(answers))
	for i := range answers {
		protos[i] = answerToProto(&answers[i])
	}

	return &quizv1.GetAllAnswersResponse{Answers: protos}, nil
}

func (h *Handler) UpdateAnswer(ctx context.Context, req *quizv1.UpdateAnswerRequest) (*quizv1.AnswerResponse, error) {
	updated, err := h.answerSvc.Update(ctx, req.GetId(), req)
	if err != nil {
		return nil, h.toAnswerGRPCError(err)
	}

	return answerToProto(updated), nil
}

func (h *Handler) DeleteAnswer(ctx context.Context, req *quizv1.DeleteAnswerRequest) (*quizv1.DeleteAnswerResponse, error) {
	if err := h.answerSvc.Delete(ctx, req.GetId()); err != nil {
		return nil, h.toAnswerGRPCError(err)
	}

	return &quizv1.DeleteAnswerResponse{Message: "answer deleted"}, nil
}

func answerToProto(a *domain.Answer) *quizv1.AnswerResponse {
	return &quizv1.AnswerResponse{
		Id:         a.ID,
		QuestionId: a.QuestionID,
		Text:       a.Text,
		Correct:    a.Correct,
		CreatedAt:  a.CreatedAt.Unix(),
		UpdatedAt:  a.UpdatedAt.Unix(),
	}
}

func (h *Handler) toAnswerGRPCError(err error) error {
	if errors.Is(err, domain.ErrAnswerNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	h.log.Error("unexpected error", zap.Error(err))
	return status.Error(codes.Internal, "internal error")
}
