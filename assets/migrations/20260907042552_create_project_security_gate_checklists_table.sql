-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gate_checklists (
    id INT IDENTITY(1,1) PRIMARY KEY,
    security_gate_id INT NOT NULL,
    requirement_id INT NOT NULL,
    response_status_id INT NULL,
    comment NVARCHAR(MAX) NULL,
    is_exception BIT NOT NULL
        CONSTRAINT DF_project_security_gate_checklists_is_exception DEFAULT 0,
    exception_justification NVARCHAR(MAX) NULL,
    created_by_id INT NULL,
    created_at DATETIME NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_project_security_gate_checklists_security_gate
        FOREIGN KEY (security_gate_id)
        REFERENCES project_security_gates(id),

    CONSTRAINT FK_project_security_gate_checklists_requirement
        FOREIGN KEY (requirement_id)
        REFERENCES project_security_gate_checklist_template_requirements(id),

    CONSTRAINT FK_project_security_gate_checklists_status
        FOREIGN KEY (response_status_id)
        REFERENCES project_security_gate_checklist_response_statuses(id),

    CONSTRAINT FK_project_security_gate_checklists_created_by
        FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_security_gate_checklists_updated_by
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT UQ_project_security_gate_checklists_gate_requirement
        UNIQUE (security_gate_id, requirement_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gate_checklists;
-- +goose StatementEnd