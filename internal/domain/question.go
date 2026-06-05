package domain

import "time"

type Question struct {
	ID             int64
	QuizID         int64
	Text           string
	Context        string
	VideoAnswerUrl string
	Order          int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type QuestionWithAnswers struct {
	Question
	Answers []Answer
}

type QuestionWithCorrectAnswer struct {
	ID                int64
	CorrectAnswerIDs  []int64
	TotalAnswersCount int64
}
