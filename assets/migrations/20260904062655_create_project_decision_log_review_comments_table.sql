-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_decision_log_review_comments (
    id INT IDENTITY(1,1) PRIMARY KEY,
    comment NVARCHAR(MAX) NOT NULL,
    is_edited BIT NOT NULL DEFAULT 0,
    project_decision_log_id INT NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_project_decision_log_review_comments_decision_logs
        FOREIGN KEY (project_decision_log_id) REFERENCES project_decision_logs(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_decision_log_review_comments_users_1
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_decision_log_review_comments_users_2
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_decision_log_review_comments;
-- +goose StatementEnd
