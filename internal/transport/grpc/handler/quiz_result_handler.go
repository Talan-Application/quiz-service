package handler

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	quizv1 "github.com/Talan-Application/proto-generation/quiz/v1"
	"github.com/Talan-Application/quiz-service/internal/domain"
	"github.com/Talan-Application/quiz-service/internal/service"
	"github.com/Talan-Application/quiz-service/internal/transport/grpc/ctxkeys"
)

// QuizResultHandler handles QuizResultService gRPC calls.
// It resolves answers from the answer service, saves the result, and returns a full evaluation.
type QuizResultHandler struct {
	quizv1.UnimplementedQuizResultServiceServer
	answerSvc    service.IAnswerService
	resultSvc    service.IQuizResultService
	log          *zap.Logger
}

func NewQuizResultHandler(
	answerSvc service.IAnswerService,
	resultSvc service.IQuizResultService,
	log *zap.Logger,
) *QuizResultHandler {
	return &QuizResultHandler{
		answerSvc: answerSvc,
		resultSvc: resultSvc,
		log:       log,
	}
}

// SubmitQuiz evaluates the submitted answers, persists the result, and returns the full breakdown.
func (h *QuizResultHandler) SubmitQuiz(ctx context.Context, req *quizv1.SubmitQuizRequest) (*quizv1.SubmitQuizResponse, error) {
	userID, ok := ctx.Value(ctxkeys.UserIDKey).(int64)
	if !ok || userID == 0 {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	if len(req.GetAnswers()) == 0 {
		return nil, status.Error(codes.InvalidArgument, "answers must not be empty")
	}

	resultAnswers := make([]domain.QuizResultAnswer, 0, len(req.GetAnswers()))
	protoResults := make([]*quizv1.QuestionResult, 0, len(req.GetAnswers()))
	correctCount := 0

	for _, sub := range req.GetAnswers() {
		allAnswers, err := h.answerSvc.GetAll(ctx, sub.GetQuestionId(), nil, nil)
		if err != nil {
			h.log.Error("failed to get answers for question",
				zap.Int64("question_id", sub.GetQuestionId()),
				zap.Error(err),
			)
			return nil, h.toResultGRPCError(err)
		}

		var correctAnswerID int64
		isCorrect := false
		for _, a := range allAnswers {
			if a.Correct {
				correctAnswerID = a.ID
			}
			if a.ID == sub.GetAnswerId() && a.Correct {
				isCorrect = true
			}
		}

		if isCorrect {
			correctCount++
		}

		resultAnswers = append(resultAnswers, domain.QuizResultAnswer{
			QuestionID:       sub.GetQuestionId(),
			SelectedAnswerID: sub.GetAnswerId(),
			CorrectAnswerID:  correctAnswerID,
			IsCorrect:        isCorrect,
		})

		protoResults = append(protoResults, &quizv1.QuestionResult{
			QuestionId:       sub.GetQuestionId(),
			SelectedAnswerId: sub.GetAnswerId(),
			CorrectAnswerId:  correctAnswerID,
			IsCorrect:        isCorrect,
		})
	}

	total := len(req.GetAnswers())
	var score float64
	if total > 0 {
		score = float64(correctCount) / float64(total) * 100
	}

	saved, err := h.resultSvc.Submit(ctx, &domain.QuizResult{
		QuizID:         req.GetQuizId(),
		UserID:         userID,
		Score:          score,
		TotalQuestions: total,
		CorrectAnswers: correctCount,
		Answers:        resultAnswers,
	})
	if err != nil {
		h.log.Error("failed to save quiz result",
			zap.Int64("quiz_id", req.GetQuizId()),
			zap.Error(err),
		)
		return nil, h.toResultGRPCError(err)
	}

	return &quizv1.SubmitQuizResponse{
		ResultId:       saved.ID,
		TotalQuestions: int32(total),
		CorrectAnswers: int32(correctCount),
		Score:          score,
		Results:        protoResults,
	}, nil
}

// GetQuizResults returns persisted result history for a quiz, optionally filtered by user.
func (h *QuizResultHandler) GetQuizResults(ctx context.Context, req *quizv1.GetQuizResultsRequest) (*quizv1.GetQuizResultsResponse, error) {
	var (
		results []domain.QuizResult
		err     error
	)

	if req.GetUserId() != 0 {
		results, err = h.resultSvc.GetByQuizAndUser(ctx, req.GetQuizId(), req.GetUserId())
	} else {
		results, err = h.resultSvc.GetByQuiz(ctx, req.GetQuizId())
	}
	if err != nil {
		return nil, h.toResultGRPCError(err)
	}

	summaries := make([]*quizv1.QuizResultSummary, len(results))
	for i, r := range results {
		summaries[i] = &quizv1.QuizResultSummary{
			Id:             r.ID,
			QuizId:         r.QuizID,
			UserId:         r.UserID,
			Score:          r.Score,
			TotalQuestions: int32(r.TotalQuestions),
			CorrectAnswers: int32(r.CorrectAnswers),
			SubmittedAt:    r.SubmittedAt.Unix(),
		}
	}

	return &quizv1.GetQuizResultsResponse{Results: summaries}, nil
}

func (h *QuizResultHandler) toResultGRPCError(err error) error {
	h.log.Error("quiz result service error", zap.Error(err))
	return status.Error(codes.Internal, "internal error")
}

