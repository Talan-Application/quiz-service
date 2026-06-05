package handler

import (
	"context"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	quizv1 "github.com/Talan-Application/proto-generation/quiz/v1"
	"github.com/Talan-Application/quiz-service/internal/domain"
	"github.com/Talan-Application/quiz-service/internal/service"
	"github.com/Talan-Application/quiz-service/internal/transport/grpc/ctxkeys"
)

type Handler struct {
	quizv1.UnimplementedQuizServiceServer
	quizv1.UnimplementedQuestionServiceServer
	quizSvc     service.IQuizService
	questionSvc service.IQuestionService
	log         *zap.Logger
}

func NewHandler(quizSvc service.IQuizService, questionSvc service.IQuestionService, log *zap.Logger) *Handler {
	return &Handler{quizSvc: quizSvc, questionSvc: questionSvc, log: log}
}

func (h *Handler) CreateQuiz(ctx context.Context, req *quizv1.CreateQuizRequest) (*quizv1.QuizResponse, error) {
	userID, ok := ctx.Value(ctxkeys.UserIDKey).(int64)
	if !ok || userID == 0 {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	created, err := h.quizSvc.Create(ctx, req, userID)
	if err != nil {
		return nil, h.toGRPCError(err)
	}

	return toProto(created), nil
}

func (h *Handler) GetQuiz(ctx context.Context, req *quizv1.GetQuizRequest) (*quizv1.QuizResponse, error) {
	quiz, err := h.quizSvc.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, h.toGRPCError(err)
	}

	return toProto(quiz), nil
}

func (h *Handler) GetAllQuizzes(ctx context.Context, req *quizv1.GetAllQuizzesRequest) (*quizv1.GetAllQuizzesResponse, error) {
	quizzes, err := h.quizSvc.GetAll(ctx, req)
	if err != nil {
		return nil, h.toGRPCError(err)
	}

	protos := make([]*quizv1.QuizResponse, len(quizzes))
	for i := range quizzes {
		protos[i] = toProto(&quizzes[i])
	}

	return &quizv1.GetAllQuizzesResponse{Quizzes: protos}, nil
}

func (h *Handler) PublishQuiz(ctx context.Context, req *quizv1.PublishQuizRequest) (*quizv1.PublishQuizResponse, error) {
	if err := h.quizSvc.Publish(ctx, req.GetId()); err != nil {
		return nil, h.toGRPCError(err)
	}

	return &quizv1.PublishQuizResponse{}, nil
}

func (h *Handler) GetMyQuizzes(ctx context.Context, req *quizv1.GetMyQuizzesRequest) (*quizv1.GetAllQuizzesResponse, error) {
	userID, ok := ctx.Value(ctxkeys.UserIDKey).(int64)
	if !ok || userID == 0 {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	quizzes, err := h.quizSvc.GetAllByAuthor(ctx, req, userID)
	if err != nil {
		return nil, h.toGRPCError(err)
	}

	protos := make([]*quizv1.QuizResponse, len(quizzes))
	for i := range quizzes {
		protos[i] = toProto(&quizzes[i])
	}

	return &quizv1.GetAllQuizzesResponse{Quizzes: protos}, nil
}

func (h *Handler) UpdateQuiz(ctx context.Context, req *quizv1.UpdateQuizRequest) (*quizv1.QuizResponse, error) {
	updated, err := h.quizSvc.Update(ctx, req.GetId(), req)
	if err != nil {
		return nil, h.toGRPCError(err)
	}

	return toProto(updated), nil
}

func (h *Handler) DeleteQuiz(ctx context.Context, req *quizv1.DeleteQuizRequest) (*quizv1.DeleteQuizResponse, error) {
	if err := h.quizSvc.Delete(ctx, req.GetId()); err != nil {
		return nil, h.toGRPCError(err)
	}

	return &quizv1.DeleteQuizResponse{Message: "quiz deleted"}, nil
}

func toProto(q *domain.Quiz) *quizv1.QuizResponse {
	return &quizv1.QuizResponse{
		Id:              q.ID,
		Title:           q.Title,
		Language:        q.Language,
		AuthorId:        q.AuthorID,
		Status:          string(q.Status),
		Type:            string(q.Type),
		CommonSubjectId: q.CommonSubjectID,
		IsEntStandard:   q.IsEntStandard,
		CreatedAt:       q.CreatedAt.Unix(),
		UpdatedAt:       q.UpdatedAt.Unix(),
	}
}

func (h *Handler) toGRPCError(err error) error {
	if errors.Is(err, domain.ErrQuizNotFound) {
		return status.Error(codes.NotFound, err.Error())
	}
	h.log.Error("unexpected error", zap.Error(err))
	return status.Error(codes.Internal, "internal error")
}
