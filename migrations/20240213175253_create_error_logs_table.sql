-- +goose Up
CREATE TABLE error_logs (
    id INT PRIMARY KEY IDENTITY (1, 1),
    level VARCHAR(45),
    URL VARCHAR(45) NULL,
    message VARCHAR(MAX),
    created_at DATETIME DEFAULT GETDATE ()
);

-- +goose Down
-- SQL statement to drop the table or undo the change in the Up migration
DROP TABLE error_logs;