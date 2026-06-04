package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Talan-Application/quiz-service/internal/domain"
)

type QuizResultRepository struct {
	db *pgxpool.Pool
}

func NewQuizResultRepository(db *pgxpool.Pool) *QuizResultRepository {
	return &QuizResultRepository{db: db}
}

func (r *QuizResultRepository) Save(ctx context.Context, result *domain.QuizResult) (*domain.QuizResult, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	const resultQuery = `
		INSERT INTO quiz_results (quiz_id, user_id, score, max_score, total_questions_count, correct_answers_count, incorrect_answers_count, unanswered_questions, submitted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
		RETURNING id, submitted_at`

	err = tx.QueryRow(ctx, resultQuery,
		result.QuizID,
		result.UserID,
		result.Score,
		result.MaxScore,
		result.TotalQuestionsCount,
		result.CorrectAnswersCount,
		result.IncorrectAnswersCount,
		result.UnansweredQuestions,
	).Scan(&result.ID, &result.SubmittedAt)
	if err != nil {
		return nil, fmt.Errorf("insert quiz result: %w", err)
	}

	const answerQuery = `
		INSERT INTO quiz_result_answers (result_id, question_id, selected_answer_ids, correct_answer_ids, score, max_score)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`

	for i := range result.Answers {
		result.Answers[i].ResultID = result.ID
		err = tx.QueryRow(ctx, answerQuery,
			result.ID,
			result.Answers[i].QuestionID,
			result.Answers[i].SelectedAnswerIDs,
			result.Answers[i].CorrectAnswerIDs,
			result.Answers[i].Score,
			result.Answers[i].MaxScore,
		).Scan(&result.Answers[i].ID)
		if err != nil {
			return nil, fmt.Errorf("insert result answer: %w", err)
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return result, nil
}

func (r *QuizResultRepository) GetByQuizAndUser(ctx context.Context, quizID, userID int64) ([]domain.QuizResult, error) {
	const query = `
		SELECT id, quiz_id, user_id, score, max_score, total_questions_count, correct_answers_count, incorrect_answers_count, unanswered_questions, submitted_at
		FROM quiz_results
		WHERE quiz_id = $1 AND user_id = $2
		ORDER BY submitted_at DESC`

	rows, err := r.db.Query(ctx, query, quizID, userID)
	if err != nil {
		return nil, fmt.Errorf("query quiz results: %w", err)
	}
	defer rows.Close()

	return scanResults(rows)
}

func (r *QuizResultRepository) GetByQuiz(ctx context.Context, quizID int64) ([]domain.QuizResult, error) {
	const query = `
		SELECT id, quiz_id, user_id, score, max_score, total_questions_count, correct_answers_count, incorrect_answers_count, unanswered_questions, submitted_at
		FROM quiz_results
		WHERE quiz_id = $1
		ORDER BY submitted_at DESC`

	rows, err := r.db.Query(ctx, query, quizID)
	if err != nil {
		return nil, fmt.Errorf("query quiz results: %w", err)
	}
	defer rows.Close()

	return scanResults(rows)
}

func scanResults(rows pgx.Rows) ([]domain.QuizResult, error) {
	var results []domain.QuizResult
	for rows.Next() {
		var r domain.QuizResult
		if err := rows.Scan(
			&r.ID,
			&r.QuizID,
			&r.UserID,
			&r.Score,
			&r.MaxScore,
			&r.TotalQuestionsCount,
			&r.CorrectAnswersCount,
			&r.IncorrectAnswersCount,
			&r.UnansweredQuestions,
			&r.SubmittedAt,
		); err != nil {
			return nil, fmt.Errorf("scan quiz result: %w", err)
		}
		results = append(results, r)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}
	return results, nil
}
