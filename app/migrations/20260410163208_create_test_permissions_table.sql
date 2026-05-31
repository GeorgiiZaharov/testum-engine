-- +goose Up
CREATE TABLE test_permissions (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    test_id INTEGER NOT NULL,
    "group" TEXT NOT NULL,

    CONSTRAINT fk_permissions_test
        FOREIGN KEY (test_id)
        REFERENCES tests(id)
        ON DELETE CASCADE
);

-- +goose Down
DROP TABLE test_permissions;
