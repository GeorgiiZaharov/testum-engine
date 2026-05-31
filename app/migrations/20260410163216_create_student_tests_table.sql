-- +goose Up
CREATE TABLE student_tests (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    student_id INTEGER NOT NULL,
    test_id INTEGER NOT NULL,

    mark INTEGER,
    "group" TEXT NOT NULL,
    success_rate REAL,

    date_start DATETIME,
    date_end DATETIME,

    CONSTRAINT fk_student_tests_user
        FOREIGN KEY (student_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_student_tests_test
        FOREIGN KEY (test_id)
        REFERENCES tests(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_student_test
        UNIQUE (student_id, test_id)
);

-- +goose Down
DROP TABLE student_tests;
