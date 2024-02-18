-- +goose Up
-- +goose StatementBegin
CREATE TABLE error_logs (
    id INT IDENTITY(1,1) PRIMARY KEY,
    http_status_code INT,
    http_request_method NVARCHAR(10),
    URL NVARCHAR(2048),
    error_message NVARCHAR(MAX),
    stack_trace NVARCHAR(MAX),
    user_id INT NULL, -- user_id can be NULL, indicating that it's optional
    created_at DATETIME2,
    -- CONSTRAINT FK_error_logs_user_id FOREIGN KEY (user_id)
    --     REFERENCES users(id) -- Assuming 'users' is the table and 'id' is the primary key column in the users table
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.error_logs;
-- +goose StatementEnd
