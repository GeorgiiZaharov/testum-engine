-- +goose Up
CREATE TABLE student_tests (
    id INT AUTO_INCREMENT PRIMARY KEY,
    student_id INT NOT NULL,
    test_id INT NOT NULL,

    mark INT NULL,
    `group` VARCHAR(100) NOT NULL,
    success_rate FLOAT NULL,

    date_start TIMESTAMP NULL,
    date_end TIMESTAMP NULL,

    CONSTRAINT fk_student_tests_user
        FOREIGN KEY (student_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_student_tests_test
        FOREIGN KEY (test_id)
        REFERENCES tests(id)
        ON DELETE CASCADE,

    CONSTRAINT uq_student_test UNIQUE (student_id, test_id)
);

-- +goose Down
DROP TABLE student_tests;
