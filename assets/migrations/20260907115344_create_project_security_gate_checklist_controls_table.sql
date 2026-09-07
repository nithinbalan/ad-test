-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gate_checklist_controls (
    id INT PRIMARY KEY IDENTITY(1, 1),
    checklist_id INT NOT NULL,
    control_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_project_security_gate_checklist_controls_checklist
        FOREIGN KEY (checklist_id)
        REFERENCES project_security_gate_checklists(id),

    CONSTRAINT FK_project_security_gate_checklist_controls_control
        FOREIGN KEY (control_id)
        REFERENCES controls(id),

    CONSTRAINT FK_project_security_gate_checklist_controls_created_by
        FOREIGN KEY (created_by_id)
        REFERENCES users(id),

    CONSTRAINT UQ_project_security_gate_checklist_controls
        UNIQUE (checklist_id, control_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gate_checklist_controls;
-- +goose StatementEnd