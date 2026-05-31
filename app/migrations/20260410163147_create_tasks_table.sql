-- +goose Up
CREATE TABLE tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    test_id INTEGER NOT NULL,
    text TEXT NOT NULL,
    image_url TEXT,
    is_hard INTEGER NOT NULL,

    CONSTRAINT fk_tasks_test
        FOREIGN KEY (test_id)
        REFERENCES tests(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE tasks;
