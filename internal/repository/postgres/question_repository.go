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

func (r *QuestionRepository) CreateWithAnswers(ctx context.Context, question *domain.Question, answers []domain.Answer) (*domain.QuestionWithAnswers, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx,
		`INSERT INTO questions (quiz_id, text, context, video_answer_url, "order", created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id, created_at, updated_at`,
		question.QuizID, question.Text, question.Context, question.VideoAnswerUrl, question.Order,
	).Scan(&question.ID, &question.CreatedAt, &question.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("create question: %w", err)
	}

	result := &domain.QuestionWithAnswers{Question: *question}

	if len(answers) > 0 {
		texts := make([]string, len(answers))
		corrects := make([]bool, len(answers))
		for i, a := range answers {
			texts[i] = a.Text
			corrects[i] = a.Correct
		}

		rows, err := tx.Query(ctx,
			`INSERT INTO answers (question_id, text, correct, created_at, updated_at)
			 SELECT $1, unnest($2::text[]), unnest($3::bool[]), NOW(), NOW()
			 RETURNING id, question_id, text, correct, created_at, updated_at`,
			question.ID, texts, corrects,
		)
		if err != nil {
			return nil, fmt.Errorf("create answers: %w", err)
		}

		for rows.Next() {
			var a domain.Answer
			if err := rows.Scan(&a.ID, &a.QuestionID, &a.Text, &a.Correct, &a.CreatedAt, &a.UpdatedAt); err != nil {
				rows.Close()
				return nil, fmt.Errorf("scan answer: %w", err)
			}
			result.Answers = append(result.Answers, a)
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("answers rows: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return result, nil
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

func (r *QuestionRepository) GetByIdWithAnswers(ctx context.Context, id int64) (*domain.QuestionWithAnswers, error) {
	question, err := r.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx,
		`SELECT id, question_id, text, correct, created_at, updated_at FROM answers WHERE question_id = $1 ORDER BY id`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("get answers for question: %w", err)
	}
	defer rows.Close()

	result := &domain.QuestionWithAnswers{Question: *question}
	for rows.Next() {
		var a domain.Answer
		if err := rows.Scan(&a.ID, &a.QuestionID, &a.Text, &a.Correct, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan answer: %w", err)
		}
		result.Answers = append(result.Answers, a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("answers rows: %w", err)
	}

	return result, nil
}

func (r *QuestionRepository) GetAllWithAnswers(ctx context.Context, quizID int64, limit *int, offset *int) ([]domain.QuestionWithAnswers, error) {
	questions, err := r.GetAll(ctx, quizID, limit, offset)
	if err != nil {
		return nil, err
	}
	if len(questions) == 0 {
		return []domain.QuestionWithAnswers{}, nil
	}

	ids := make([]int64, len(questions))
	for i, q := range questions {
		ids[i] = q.ID
	}

	rows, err := r.db.Query(ctx,
		`SELECT id, question_id, text, correct, created_at, updated_at FROM answers WHERE question_id = ANY($1) ORDER BY question_id, id`,
		ids,
	)
	if err != nil {
		return nil, fmt.Errorf("get answers for questions: %w", err)
	}
	defer rows.Close()

	answerMap := make(map[int64][]domain.Answer)
	for rows.Next() {
		var a domain.Answer
		if err := rows.Scan(&a.ID, &a.QuestionID, &a.Text, &a.Correct, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan answer: %w", err)
		}
		answerMap[a.QuestionID] = append(answerMap[a.QuestionID], a)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("answers rows: %w", err)
	}

	results := make([]domain.QuestionWithAnswers, len(questions))
	for i, q := range questions {
		results[i] = domain.QuestionWithAnswers{
			Question: q,
			Answers:  answerMap[q.ID],
		}
	}

	return results, nil
}

func (r *QuestionRepository) UpdateWithAnswers(ctx context.Context, id int64, question *domain.Question, answers []domain.Answer) (*domain.QuestionWithAnswers, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	err = tx.QueryRow(ctx,
		`UPDATE questions SET text = $1, context = $2, video_answer_url = $3, "order" = $4, updated_at = NOW()
		 WHERE id = $5 RETURNING id, quiz_id, text, context, video_answer_url, "order", created_at, updated_at`,
		question.Text, question.Context, question.VideoAnswerUrl, question.Order, id,
	).Scan(&question.ID, &question.QuizID, &question.Text, &question.Context, &question.VideoAnswerUrl, &question.Order, &question.CreatedAt, &question.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("update question: %w", domain.ErrQuestionNotFound)
		}
		return nil, fmt.Errorf("update question: %w", err)
	}

	if _, err = tx.Exec(ctx, `DELETE FROM answers WHERE question_id = $1`, id); err != nil {
		return nil, fmt.Errorf("delete answers: %w", err)
	}

	result := &domain.QuestionWithAnswers{Question: *question}

	if len(answers) > 0 {
		texts := make([]string, len(answers))
		corrects := make([]bool, len(answers))
		for i, a := range answers {
			texts[i] = a.Text
			corrects[i] = a.Correct
		}

		answerRows, err := tx.Query(ctx,
			`INSERT INTO answers (question_id, text, correct, created_at, updated_at)
			 SELECT $1, unnest($2::text[]), unnest($3::bool[]), NOW(), NOW()
			 RETURNING id, question_id, text, correct, created_at, updated_at`,
			id, texts, corrects,
		)
		if err != nil {
			return nil, fmt.Errorf("create answers: %w", err)
		}

		for answerRows.Next() {
			var a domain.Answer
			if err := answerRows.Scan(&a.ID, &a.QuestionID, &a.Text, &a.Correct, &a.CreatedAt, &a.UpdatedAt); err != nil {
				answerRows.Close()
				return nil, fmt.Errorf("scan answer: %w", err)
			}
			result.Answers = append(result.Answers, a)
		}
		answerRows.Close()
		if err := answerRows.Err(); err != nil {
			return nil, fmt.Errorf("answers rows: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return result, nil
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

func (r *QuestionRepository) GetByQuizId(ctx context.Context, quizId int64) ([]domain.Question, error) {
	query := `SELECT id, quiz_id, text, context, video_answer_url, "order", created_at, updated_at
			  FROM questions WHERE quiz_id = $1 ORDER BY "order", id`

	rows, err := r.db.Query(ctx, query, quizId)
	if err != nil {
		return nil, fmt.Errorf("get questions by quiz id: %w", err)
	}
	defer rows.Close()

	var questions []domain.Question
	for rows.Next() {
		var q domain.Question
		if err := rows.Scan(&q.ID, &q.QuizID, &q.Text, &q.Context, &q.VideoAnswerUrl, &q.Order, &q.CreatedAt, &q.UpdatedAt); err != nil {
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
