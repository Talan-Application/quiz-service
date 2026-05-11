CREATE TABLE IF NOT EXISTS questions
(
    id               BIGSERIAL    PRIMARY KEY,
    quiz_id          BIGINT       NOT NULL REFERENCES quizzes(id) ON DELETE CASCADE,
    text             TEXT         NOT NULL,
    context          TEXT,
    video_answer_url VARCHAR(500),
    "order"          INT          NOT NULL DEFAULT 0,
    created_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_questions_quiz_id ON questions (quiz_id);
CREATE INDEX idx_questions_order   ON questions (quiz_id, "order");
