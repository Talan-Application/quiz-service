package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Talan-Application/quiz-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuestionRepository struct {
	db *pgxpool.Pool
}

func NewQuestionRepository(db *pgxpool.Pool) *QuestionRepository {
	return &QuestionRepository{db}
}

func (r *QuestionRepository) Create(ctx context.Context, question *domain.Question) (*domain.Question, error) {
	query := `INSERT INTO questions (quiz_id, text, context, video_answer_url, "order", created_at, updated_at)
			  VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		question.QuizID,
		question.Text,
		question.Context,
		question.VideoAnswerUrl,
		question.Order,
	).Scan(&question.ID, &question.CreatedAt, &question.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("create question: %w", err)
	}

	return question, nil
}

func (r *QuestionRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM questions WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete question: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("question with id=%d: %w", id, domain.ErrQuestionNotFound)
	}

	return nil
}

func (r *QuestionRepository) Update(ctx context.Context, id int64, question *domain.Question) (*domain.Question, error) {
	query := `UPDATE questions SET text = $1, context = $2, video_answer_url = $3, "order" = $4,
			  updated_at = NOW() WHERE id = $5 RETURNING updated_at`

	err := r.db.QueryRow(ctx, query,
		question.Text,
		question.Context,
		question.VideoAnswerUrl,
		question.Order,
		id,
	).Scan(&question.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("update question: %w", domain.ErrQuestionNotFound)
		}
		return nil, fmt.Errorf("update question: %w", err)
	}

	question.ID = id
	return question, nil
}

func (r *QuestionRepository) GetAll(ctx context.Context, quizID int64, limit *int, offset *int) ([]domain.Question, error) {
	query := `SELECT id, quiz_id, text, context, video_answer_url, "order", created_at, updated_at
			  FROM questions WHERE quiz_id = $1 ORDER BY "order", id LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, quizID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select questions: %w", err)
	}
	defer rows.Close()

	var questions []domain.Question
	for rows.Next() {
		var q domain.Question
		err := rows.Scan(&q.ID, &q.QuizID, &q.Text, &q.Context, &q.VideoAnswerUrl, &q.Order, &q.CreatedAt, &q.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan question: %w", err)
		}
		questions = append(questions, q)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return questions, nil
}

func (r *QuestionRepository) GetById(ctx context.Context, id int64) (*domain.Question, error) {
	query := `SELECT id, quiz_id, text, context, video_answer_url, "order", created_at, updated_at
			  FROM questions WHERE id = $1`

	q := &domain.Question{}
	err := r.db.QueryRow(ctx, query, id).
		Scan(&q.ID, &q.QuizID, &q.Text, &q.Context, &q.VideoAnswerUrl, &q.Order, &q.CreatedAt, &q.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get question: %w", domain.ErrQuestionNotFound)
		}
		return nil, fmt.Errorf("get question: %w", err)
	}

	return q, nil
}
