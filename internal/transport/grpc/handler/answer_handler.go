package handler

import (
	quizv1 "github.com/Talan-Application/proto-generation/quiz/v1"
	"github.com/Talan-Application/quiz-service/internal/domain"
)

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
