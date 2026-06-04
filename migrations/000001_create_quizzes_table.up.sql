CREATE TYPE quiz_status AS ENUM ('draft', 'published');
CREATE TYPE quiz_type AS ENUM ('ent', 'monthly_exam', 'exam');

CREATE TABLE IF NOT EXISTS quizzes
(
    id                BIGSERIAL PRIMARY KEY,
    title             VARCHAR(255) NOT NULL,
    language          VARCHAR(50)  NOT NULL,
    author_id         BIGINT       NOT NULL,
    status            quiz_status  NOT NULL DEFAULT 'draft',
    type              quiz_type    NOT NULL,
    common_subject_id BIGINT       NOT NULL,
    is_ent_standard   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_quizzes_author_id ON quizzes (author_id);
CREATE INDEX idx_quizzes_subject_id ON quizzes (common_subject_id);
CREATE INDEX idx_quizzes_status ON quizzes (status);
