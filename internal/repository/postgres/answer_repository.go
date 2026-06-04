package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Talan-Application/quiz-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AnswerRepository struct {
	db *pgxpool.Pool
}

func NewAnswerRepository(db *pgxpool.Pool) *AnswerRepository {
	return &AnswerRepository{db}
}

func (r *AnswerRepository) Create(ctx context.Context, answer *domain.Answer) (*domain.Answer, error) {
	query := `INSERT INTO answers (question_id, text, correct, created_at, updated_at)
			  VALUES ($1, $2, $3, NOW(), NOW()) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		answer.QuestionID,
		answer.Text,
		answer.Correct,
	).Scan(&answer.ID, &answer.CreatedAt, &answer.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("create answer: %w", err)
	}

	return answer, nil
}

func (r *AnswerRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM answers WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete answer: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("answer with id=%d: %w", id, domain.ErrAnswerNotFound)
	}

	return nil
}

func (r *AnswerRepository) Update(ctx context.Context, id int64, answer *domain.Answer) (*domain.Answer, error) {
	query := `UPDATE answers SET text = $1, correct = $2, updated_at = NOW()
			  WHERE id = $3 RETURNING updated_at`

	err := r.db.QueryRow(ctx, query,
		answer.Text,
		answer.Correct,
		id,
	).Scan(&answer.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("update answer: %w", domain.ErrAnswerNotFound)
		}
		return nil, fmt.Errorf("update answer: %w", err)
	}

	answer.ID = id
	return answer, nil
}

func (r *AnswerRepository) GetAll(ctx context.Context, questionID int64, limit *int, offset *int) ([]domain.Answer, error) {
	query := `SELECT id, question_id, text, correct, created_at, updated_at
			  FROM answers WHERE question_id = $1 ORDER BY id LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, questionID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select answers: %w", err)
	}
	defer rows.Close()

	var answers []domain.Answer
	for rows.Next() {
		var a domain.Answer
		err := rows.Scan(&a.ID, &a.QuestionID, &a.Text, &a.Correct, &a.CreatedAt, &a.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan answer: %w", err)
		}
		answers = append(answers, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return answers, nil
}

func (r *AnswerRepository) GetByQuestionIDs(ctx context.Context, questionIDs []int64) ([]domain.Answer, error) {
	query := `SELECT id, question_id, text, correct, created_at, updated_at
			  FROM answers WHERE question_id = ANY($1) ORDER BY question_id, id`

	rows, err := r.db.Query(ctx, query, questionIDs)
	if err != nil {
		return nil, fmt.Errorf("get answers by question ids: %w", err)
	}
	defer rows.Close()

	var answers []domain.Answer
	for rows.Next() {
		var a domain.Answer
		if err := rows.Scan(&a.ID, &a.QuestionID, &a.Text, &a.Correct, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan answer: %w", err)
		}
		answers = append(answers, a)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return answers, nil
}

func (r *AnswerRepository) GetById(ctx context.Context, id int64) (*domain.Answer, error) {
	query := `SELECT id, question_id, text, correct, created_at, updated_at
			  FROM answers WHERE id = $1`

	a := &domain.Answer{}
	err := r.db.QueryRow(ctx, query, id).
		Scan(&a.ID, &a.QuestionID, &a.Text, &a.Correct, &a.CreatedAt, &a.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get answer: %w", domain.ErrAnswerNotFound)
		}
		return nil, fmt.Errorf("get answer: %w", err)
	}

	return a, nil
}
