-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gate_checklist_exception_workflows (
    id INT IDENTITY(1,1) PRIMARY KEY,
    checklist_id INT NOT NULL,
    level INT NOT NULL,
    workflow_level_type_id INT NOT NULL,
    approval_types NVARCHAR(50) NOT NULL,
    workflow_status_id INT NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_psg_checklist_exception_workflows_checklists
        FOREIGN KEY (checklist_id) REFERENCES project_security_gate_checklists(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_psg_checklist_exception_workflows_workflow_level_types
        FOREIGN KEY (workflow_level_type_id) REFERENCES workflow_level_types(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_psg_checklist_exception_workflows_workflow_statuses
        FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_psg_checklist_exception_workflows_created_by
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_psg_checklist_exception_workflows_updated_by
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gate_checklist_exception_workflows;
-- +goose StatementEnd
