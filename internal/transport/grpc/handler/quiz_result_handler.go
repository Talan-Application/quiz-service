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

type QuizResultHandler struct {
	quizv1.UnimplementedQuizResultServiceServer
	resultSvc service.IQuizResultService
	log       *zap.Logger
}

func NewQuizResultHandler(
	resultSvc service.IQuizResultService,
	log *zap.Logger,
) *QuizResultHandler {
	return &QuizResultHandler{
		resultSvc: resultSvc,
		log:       log,
	}
}

func (h *QuizResultHandler) SubmitQuiz(ctx context.Context, req *quizv1.SubmitQuizRequest) (*quizv1.SubmitQuizResponse, error) {
	userID, ok := ctx.Value(ctxkeys.UserIDKey).(int64)
	if !ok || userID == 0 {
		return nil, status.Error(codes.Unauthenticated, "authentication required")
	}

	result, err := h.resultSvc.Submit(ctx, userID, req)
	if err != nil {
		return nil, h.toResultGRPCError(err)
	}

	questionResults := make([]*quizv1.QuestionResult, len(result.Answers))
	for i, a := range result.Answers {
		questionResults[i] = &quizv1.QuestionResult{
			QuestionId:        a.QuestionID,
			SelectedAnswerIds: a.SelectedAnswerIDs,
			CorrectAnswerIds:  a.CorrectAnswerIDs,
			Score:             a.Score,
			MaxScore:          a.MaxScore,
		}
	}

	return &quizv1.SubmitQuizResponse{
		ResultId:              result.ID,
		TotalQuestionsCount:   int32(result.TotalQuestionsCount),
		CorrectAnswersCount:   int32(result.CorrectAnswersCount),
		IncorrectAnswersCount: int32(result.IncorrectAnswersCount),
		UnansweredQuestions:   int32(result.UnansweredQuestions),
		Score:                 result.Score,
		MaxScore:              result.MaxScore,
		Results:               questionResults,
	}, nil
}

func (h *QuizResultHandler) GetQuizResults(ctx context.Context, req *quizv1.GetQuizResultsRequest) (*quizv1.GetQuizResultsResponse, error) {
	results, err := h.resultSvc.GetResults(ctx, req.GetQuizId(), req.GetUserId())
	if err != nil {
		return nil, h.toResultGRPCError(err)
	}

	summaries := make([]*quizv1.QuizResultSummary, len(results))
	for i, r := range results {
		summaries[i] = &quizv1.QuizResultSummary{
			Id:                    r.ID,
			QuizId:                r.QuizID,
			UserId:                r.UserID,
			Score:                 r.Score,
			MaxScore:              r.MaxScore,
			TotalQuestionsCount:   int32(r.TotalQuestionsCount),
			CorrectAnswersCount:   int32(r.CorrectAnswersCount),
			IncorrectAnswersCount: int32(r.IncorrectAnswersCount),
			UnansweredQuestions:   int32(r.UnansweredQuestions),
			SubmittedAt:           r.SubmittedAt.Unix(),
		}
	}

	return &quizv1.GetQuizResultsResponse{Results: summaries}, nil
}

func (h *QuizResultHandler) toResultGRPCError(err error) error {
	if errors.Is(err, domain.ErrEmptyAnswers) {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	h.log.Error("quiz result service error", zap.Error(err))
	return status.Error(codes.Internal, "internal error")
}
