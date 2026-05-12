package domain

import "time"

type QuizResult struct {
	ID             int64
	QuizID         int64
	UserID         int64
	Score          float64
	TotalQuestions int
	CorrectAnswers int
	SubmittedAt    time.Time
	Answers        []QuizResultAnswer
}

type QuizResultAnswer struct {
	ID               int64
	ResultID         int64
	QuestionID       int64
	SelectedAnswerID int64
	CorrectAnswerID  int64
	IsCorrect        bool
}
