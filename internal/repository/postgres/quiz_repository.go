package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Talan-Application/quiz-service/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type QuizRepository struct {
	db *pgxpool.Pool
}

func NewQuizRepository(db *pgxpool.Pool) *QuizRepository {
	return &QuizRepository{db}
}

func (r *QuizRepository) Create(ctx context.Context, quiz *domain.Quiz) (*domain.Quiz, error) {
	query := `INSERT INTO quizzes (title, language, author_id, type, common_subject_id, is_ent_standard, created_at, updated_at)
				VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		ctx,
		query,
		quiz.Title,
		quiz.Language,
		quiz.AuthorID,
		quiz.Type,
		quiz.CommonSubjectID,
		quiz.IsEntStandard,
	).Scan(&quiz.ID, &quiz.CreatedAt, &quiz.UpdatedAt)

	if err != nil {
		if isUniqueViolation(err) {
			return nil, fmt.Errorf("repository: %w", domain.ErrQuizNotFound)
		}
		return nil, fmt.Errorf("get from db: %w", err)
	}

	return quiz, nil
}

func (r *QuizRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM quizzes WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)

	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("user with id=%d: %w", id, domain.ErrQuizNotFound)
	}

	return nil
}

func (r *QuizRepository) Update(ctx context.Context, id int64, quiz *domain.Quiz) (*domain.Quiz, error) {
	query := `UPDATE quizzes SET title = $1, language = $2, type = $3, is_ent_standard = $4,
	 		updated_at = NOW() WHERE id = $5 RETURNING updated_at`

	err := r.db.QueryRow(ctx, query, quiz.Title, quiz.Language, quiz.Type, quiz.IsEntStandard, id).
		Scan(&quiz.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("update quiz: %w", domain.ErrQuizNotFound)
		}
		return nil, fmt.Errorf("update quiz: %w", err)
	}

	return quiz, nil
}

func (r *QuizRepository) GetAll(ctx context.Context, status *domain.QuizStatus, limit *int, offset *int) ([]domain.Quiz, error) {
	var (
		rows pgx.Rows
		err  error
	)

	if status == nil {
		rows, err = r.db.Query(ctx,
			`SELECT id, title, language, author_id, status, type, common_subject_id, is_ent_standard, created_at, updated_at
			 FROM quizzes ORDER BY id LIMIT $1 OFFSET $2`,
			limit, offset,
		)
	} else {
		rows, err = r.db.Query(ctx,
			`SELECT id, title, language, author_id, status, type, common_subject_id, is_ent_standard, created_at, updated_at
			 FROM quizzes WHERE status = $1::quiz_status ORDER BY id LIMIT $2 OFFSET $3`,
			string(*status), limit, offset,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("select quiz: %w", err)
	}
	defer rows.Close()

	var quizzes []domain.Quiz
	for rows.Next() {
		var quiz domain.Quiz
		if err := rows.Scan(&quiz.ID, &quiz.Title, &quiz.Language, &quiz.AuthorID,
			&quiz.Status, &quiz.Type, &quiz.CommonSubjectID, &quiz.IsEntStandard, &quiz.CreatedAt, &quiz.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan quiz: %w", err)
		}
		quizzes = append(quizzes, quiz)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return quizzes, nil
}

func (r *QuizRepository) Publish(ctx context.Context, id int64) error {
	query := `UPDATE quizzes SET status = 'published', updated_at = NOW() WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("publish quiz: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("publish quiz: %w", domain.ErrQuizNotFound)
	}

	return nil
}

func (r *QuizRepository) GetAllByAuthor(ctx context.Context, authorID int64, limit *int, offset *int) ([]domain.Quiz, error) {
	query := `SELECT id, title, language, author_id, status, type, common_subject_id, is_ent_standard, created_at, updated_at
			FROM quizzes
			WHERE author_id = $1
			ORDER BY id LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, authorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("select quiz by author: %w", err)
	}
	defer rows.Close()

	var quizzes []domain.Quiz
	for rows.Next() {
		var quiz domain.Quiz
		if err := rows.Scan(&quiz.ID, &quiz.Title, &quiz.Language, &quiz.AuthorID,
			&quiz.Status, &quiz.Type, &quiz.CommonSubjectID, &quiz.IsEntStandard, &quiz.CreatedAt, &quiz.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan quiz: %w", err)
		}
		quizzes = append(quizzes, quiz)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("next rows: %w", err)
	}

	return quizzes, nil
}

func (r *QuizRepository) GetById(ctx context.Context, id int64) (*domain.Quiz, error) {
	query := `SELECT id, title, language, author_id, status, type, common_subject_id, is_ent_standard, created_at, updated_at
			FROM quizzes WHERE id = $1`

	quiz := &domain.Quiz{}
	err := r.db.QueryRow(ctx, query, id).
		Scan(&quiz.ID, &quiz.Title, &quiz.Language, &quiz.AuthorID, &quiz.Status,
			&quiz.Type, &quiz.CommonSubjectID, &quiz.IsEntStandard, &quiz.CreatedAt, &quiz.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get quiz: %w", domain.ErrQuizNotFound)
		}
		return nil, fmt.Errorf("get quiz: %w", err)
	}

	return quiz, nil
}

func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == "23505"
}
