-- +goose Up
-- +goose StatementBegin
CREATE TABLE violation_report_action_plan_comment_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    violation_report_action_plan_comment_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    -- Foreign Key Constraints
    CONSTRAINT FK_violation_report_action_plan_comment_documents FOREIGN KEY (violation_report_action_plan_comment_id) 
        REFERENCES violation_report_action_plan_comments(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_report_action_plan_comment_documents_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_violation_report_action_plan_comment_documents_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE violation_report_action_plan_comment_documents;
-- +goose StatementEnd
