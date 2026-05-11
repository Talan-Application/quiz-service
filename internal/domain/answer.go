package domain

import "time"

type Answer struct {
	ID         int64
	QuestionID int64
	Text       string
	Correct    bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
