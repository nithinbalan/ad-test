-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gate_checklist_template_requirements (
    id INT PRIMARY KEY IDENTITY(1, 1),
    template_version_id INT NOT NULL,
    title_en NVARCHAR(500) NOT NULL,
    title_ar NVARCHAR(500) NOT NULL,
    is_critical BIT NOT NULL
        CONSTRAINT DF_project_security_gate_checklist_template_requirements_is_critical DEFAULT 0,

    CONSTRAINT FK_project_security_gate_checklist_template_requirements_version
        FOREIGN KEY (template_version_id)
        REFERENCES project_security_gate_checklist_template_versions(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gate_checklist_template_requirements;
-- +goose StatementEnd
