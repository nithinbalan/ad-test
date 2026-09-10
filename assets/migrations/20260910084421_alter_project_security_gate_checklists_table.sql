-- +goose Up
-- +goose StatementBegin
ALTER TABLE project_security_gate_checklists
ADD
    submitted_by_id INT NULL,
    submitted_at DATETIME NULL,
    next_review_level INT NULL,
    workflow_status_id INT NULL,
    approved_by_id INT NULL,
    approved_at DATETIME NULL,

    CONSTRAINT FK_project_security_gate_checklists_submitted_by
        FOREIGN KEY (submitted_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_security_gate_checklists_workflow_status
        FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_security_gate_checklists_approved_by
        FOREIGN KEY (approved_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE project_security_gate_checklists
DROP CONSTRAINT FK_project_security_gate_checklists_submitted_by;

ALTER TABLE project_security_gate_checklists
DROP CONSTRAINT FK_project_security_gate_checklists_workflow_status;

ALTER TABLE project_security_gate_checklists
DROP CONSTRAINT FK_project_security_gate_checklists_approved_by;

ALTER TABLE project_security_gate_checklists
DROP COLUMN submitted_by_id,
            submitted_at,
            next_review_level,
            workflow_status_id,
            approved_by_id,
            approved_at;
-- +goose StatementEnd
