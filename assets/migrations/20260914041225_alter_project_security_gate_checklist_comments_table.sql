-- +goose Up
-- +goose StatementBegin
ALTER TABLE project_security_gate_checklist_comments
ADD
    workflow_history_id INT NULL,

    CONSTRAINT FK_psg_checklist_comments_exception_workflow_history
        FOREIGN KEY (workflow_history_id) REFERENCES project_security_gate_checklist_exception_workflow_histories(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE project_security_gate_checklist_comments
DROP CONSTRAINT FK_psg_checklist_comments_exception_workflow_history;

ALTER TABLE project_security_gate_checklist_comments
DROP COLUMN workflow_history_id;
-- +goose StatementEnd
