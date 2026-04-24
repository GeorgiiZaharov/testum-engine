-- +goose Up
CREATE TABLE tasks (
    id INT AUTO_INCREMENT PRIMARY KEY,
    test_id INT NOT NULL,
    text TEXT NOT NULL,
    image_url VARCHAR(255) NULL,
    is_hard BOOLEAN NOT NULL,

    CONSTRAINT fk_tasks_test
        FOREIGN KEY (test_id)
        REFERENCES tests(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE tasks;
