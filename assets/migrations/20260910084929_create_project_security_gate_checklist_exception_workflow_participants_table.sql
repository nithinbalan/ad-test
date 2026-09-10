-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gate_checklist_exception_workflow_participants (
    id INT IDENTITY(1,1) PRIMARY KEY,
    exception_workflow_id INT NOT NULL,
    role_id INT NULL,
    user_id INT NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_psg_checklist_exception_workflow_participants_workflows
        FOREIGN KEY (exception_workflow_id) REFERENCES project_security_gate_checklist_exception_workflows(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_psg_checklist_exception_workflow_participants_roles
        FOREIGN KEY (role_id) REFERENCES roles(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_psg_checklist_exception_workflow_participants_user
        FOREIGN KEY (user_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_psg_checklist_exception_workflow_participants_created_by
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_psg_checklist_exception_workflow_participants_updated_by
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gate_checklist_exception_workflow_participants;
-- +goose StatementEnd
