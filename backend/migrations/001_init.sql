CREATE TABLE IF NOT EXISTS todos (
    id         BIGSERIAL PRIMARY KEY,
    title      TEXT        NOT NULL,
    done       BOOLEAN     NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO todos (title, done) VALUES
    ('Buy groceries', false),
    ('Walk the dog', true),
    ('Read a book', false)
ON CONFLICT DO NOTHING;
