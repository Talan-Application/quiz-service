CREATE TABLE IF NOT EXISTS quiz_results (
    id                      BIGSERIAL PRIMARY KEY,
    quiz_id                 BIGINT        NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    user_id                 BIGINT        NOT NULL,
    score                   DECIMAL(5, 2) NOT NULL,
    max_score               DECIMAL(5, 2) NOT NULL DEFAULT 0,
    total_questions_count   INT           NOT NULL,
    correct_answers_count   INT           NOT NULL,
    incorrect_answers_count INT           NOT NULL DEFAULT 0,
    unanswered_questions    INT           NOT NULL DEFAULT 0,
    submitted_at            TIMESTAMPTZ   NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS quiz_result_answers (
    id                  BIGSERIAL PRIMARY KEY,
    result_id           BIGINT        NOT NULL REFERENCES quiz_results(id) ON DELETE CASCADE,
    question_id         BIGINT        NOT NULL,
    selected_answer_ids BIGINT[]      NOT NULL DEFAULT '{}',
    correct_answer_ids  BIGINT[]      NOT NULL DEFAULT '{}',
    score               DECIMAL(5, 2) NOT NULL DEFAULT 0,
    max_score           DECIMAL(5, 2) NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_quiz_results_quiz_id          ON quiz_results(quiz_id);
CREATE INDEX IF NOT EXISTS idx_quiz_results_user_id          ON quiz_results(user_id);
CREATE INDEX IF NOT EXISTS idx_quiz_result_answers_result_id ON quiz_result_answers(result_id);
