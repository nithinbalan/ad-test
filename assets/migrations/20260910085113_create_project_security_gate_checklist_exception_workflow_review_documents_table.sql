-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gate_checklist_exception_workflow_review_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    history_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_psg_checklist_exception_review_documents_history
        FOREIGN KEY (history_id)
        REFERENCES project_security_gate_checklist_exception_workflow_histories(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_psg_checklist_exception_review_documents_created_by
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gate_checklist_exception_workflow_review_documents;
-- +goose StatementEnd
