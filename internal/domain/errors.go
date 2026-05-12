package domain

import "errors"

var (
	ErrQuizNotFound       = errors.New("quiz not found")
	ErrQuestionNotFound   = errors.New("question not found")
	ErrAnswerNotFound     = errors.New("answer not found")
	ErrQuizResultNotFound = errors.New("quiz result not found")
)
