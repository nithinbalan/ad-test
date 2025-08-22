-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_delegation_history (
    id INT PRIMARY KEY IDENTITY(1, 1),
    user_id INT NOT NULL,
    delegatee_user_id INT NOT NULL,
    comments NVARCHAR(MAX) NULL,
    assigned_date DATE NOT NULL,
    module_count INT NOT NULL,
    task_count INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_user_delegation_logs_users_1 FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_user_delegation_logs_users_2 FOREIGN KEY (delegatee_user_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_user_delegation_logs_users_3 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_user_delegation_logs_users_4 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE user_delegation_logs;
-- +goose StatementEnd
