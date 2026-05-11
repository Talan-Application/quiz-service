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
