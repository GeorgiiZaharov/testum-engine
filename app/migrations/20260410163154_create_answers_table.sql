-- +goose Up
CREATE TABLE answers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    task_id INTEGER NOT NULL,
    text TEXT NOT NULL,
    image_url TEXT,
    is_correct INTEGER NOT NULL,

    CONSTRAINT fk_answers_task
        FOREIGN KEY (task_id)
        REFERENCES tasks(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE answers;
