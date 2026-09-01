-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gate_checklist_templates (
    id INT PRIMARY KEY IDENTITY(1, 1),
    stage_id INT NOT NULL,
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,

    CONSTRAINT UQ_project_security_gate_checklist_templates_stage
        UNIQUE (stage_id),

    CONSTRAINT FK_project_security_gate_checklist_templates_stage
        FOREIGN KEY (stage_id)
        REFERENCES project_security_gate_stages(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gate_checklist_templates;
-- +goose StatementEnd
