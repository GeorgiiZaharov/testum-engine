-- +goose Up
CREATE TABLE student_answers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    student_id INTEGER NOT NULL,
    answer_id INTEGER NOT NULL,

    date_created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_student_answers_user
        FOREIGN KEY (student_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_student_answers_answer
        FOREIGN KEY (answer_id)
        REFERENCES answers(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE student_answers;
