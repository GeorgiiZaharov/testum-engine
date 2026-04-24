-- +goose Up
CREATE TABLE answers (
    id INT AUTO_INCREMENT PRIMARY KEY,
    task_id INT NOT NULL,
    text TEXT NOT NULL,
    image_url VARCHAR(255) NULL,
    is_correct BOOLEAN NOT NULL,

    CONSTRAINT fk_answers_task
        FOREIGN KEY (task_id)
        REFERENCES tasks(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE answers;
