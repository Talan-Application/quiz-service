package domain

import "time"

type QuizResult struct {
	ID                    int64
	QuizID                int64
	UserID                int64
	Score                 float64
	MaxScore              float64
	TotalQuestionsCount   int
	CorrectAnswersCount   int
	IncorrectAnswersCount int
	UnansweredQuestions   int
	SubmittedAt           time.Time
	Answers               []QuizResultAnswer
}

type QuizResultAnswer struct {
	ID                int64
	ResultID          int64
	QuestionID        int64
	SelectedAnswerIDs []int64
	CorrectAnswerIDs  []int64
	Score             float64
	MaxScore          float64
}
