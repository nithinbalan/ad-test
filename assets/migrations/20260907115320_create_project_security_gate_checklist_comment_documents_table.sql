-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gate_checklist_comment_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    project_security_gate_checklist_comment_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_project_security_gate_checklist_comment_documents_comment
        FOREIGN KEY (project_security_gate_checklist_comment_id)
        REFERENCES project_security_gate_checklist_comments(id),

    CONSTRAINT FK_project_security_gate_checklist_comment_documents_created_by
        FOREIGN KEY (created_by_id)
        REFERENCES users(id),

    CONSTRAINT FK_project_security_gate_checklist_comment_documents_updated_by
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gate_checklist_comment_documents;
-- +goose StatementEnd