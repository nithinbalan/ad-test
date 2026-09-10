-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gate_checklist_exception_workflow_histories (
    id INT IDENTITY(1,1) PRIMARY KEY,
    checklist_id INT NOT NULL,
    workflow_status_id INT NOT NULL,
    level INT NOT NULL,
    comment NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),

    CONSTRAINT FK_psg_checklist_exception_workflow_histories_checklists
        FOREIGN KEY (checklist_id) REFERENCES project_security_gate_checklists(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_psg_checklist_exception_workflow_histories_workflow_statuses
        FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_psg_checklist_exception_workflow_histories_created_by
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gate_checklist_exception_workflow_histories;
-- +goose StatementEnd
