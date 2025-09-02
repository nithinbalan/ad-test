-- +goose Up
-- +goose StatementBegin
CREATE TABLE chat_history (
    id INT PRIMARY KEY IDENTITY(1, 1),
    context_id NVARCHAR(255) NOT NULL,
    chat_type NVARCHAR(100) NOT NULL,
    question NVARCHAR(MAX) NOT NULL,
    chat_query NVARCHAR(MAX) NULL,
    chat_answer NVARCHAR(MAX) NULL,
    chat_context NVARCHAR(MAX) NULL,
    chat_summary NVARCHAR(MAX) NULL,
    chat_status INT NULL,
    user_id INT NULL,
    primary_chat BIT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_chat_history_users1 FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE chat_history;
-- +goose StatementEnd