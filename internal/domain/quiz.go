package domain

import "time"

type Quiz struct {
	ID        int64
	Title     string
	Language  string
	AuthorID  int64
	Status    QuizStatus
	Type      QuizType
	SubjectID int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

type QuizStatus string

const (
	QuizStatusDraft     QuizStatus = "draft"
	QuizStatusPublished QuizStatus = "published"
)

type QuizType string

const (
	QuizTypeEnt         QuizType = "ent"
	QuizTypeMonthlyExam QuizType = "monthly_exam"
	QuizTypeExam        QuizType = "exam"
)
