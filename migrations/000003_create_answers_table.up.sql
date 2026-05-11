CREATE TABLE IF NOT EXISTS answers
(
    id          BIGSERIAL   PRIMARY KEY,
    question_id BIGINT      NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
    text        TEXT        NOT NULL,
    correct     BOOLEAN     NOT NULL DEFAULT false,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_answers_question_id ON answers (question_id);
