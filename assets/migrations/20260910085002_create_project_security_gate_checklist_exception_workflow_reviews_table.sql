-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gate_checklist_exception_workflow_reviews (
    id INT IDENTITY(1,1) PRIMARY KEY,
    exception_workflow_id INT NOT NULL,
    workflow_status_id INT NOT NULL,
    comment NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_psg_checklist_exception_workflow_reviews_workflows
        FOREIGN KEY (exception_workflow_id) REFERENCES project_security_gate_checklist_exception_workflows(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_psg_checklist_exception_workflow_reviews_workflow_statuses
        FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_psg_checklist_exception_workflow_reviews_created_by
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_psg_checklist_exception_workflow_reviews_updated_by
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gate_checklist_exception_workflow_reviews;
-- +goose StatementEnd
