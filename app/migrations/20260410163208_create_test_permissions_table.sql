-- +goose Up
CREATE TABLE test_permissions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    test_id INT NOT NULL,
    `group` VARCHAR(100) NOT NULL,

    CONSTRAINT fk_permissions_test
        FOREIGN KEY (test_id)
        REFERENCES tests(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE test_permissions;
