-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gate_checklist_template_versions (
    id INT PRIMARY KEY IDENTITY(1, 1),
    template_id INT NOT NULL,
    version INT NOT NULL,
    status VARCHAR(10) NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    published_by_id INT NULL,
    published_at DATETIME NULL,

    CONSTRAINT UQ_project_security_gate_checklist_template_versions
        UNIQUE (template_id, version),

    CONSTRAINT CHK_project_security_gate_checklist_template_versions_status
        CHECK (status IN ('draft', 'published')),

    CONSTRAINT FK_project_security_gate_checklist_template_versions_template
        FOREIGN KEY (template_id)
        REFERENCES project_security_gate_checklist_templates(id),

    CONSTRAINT FK_project_security_gate_checklist_template_versions_users_1
        FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_security_gate_checklist_template_versions_users_2
        FOREIGN KEY (published_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gate_checklist_template_versions;
-- +goose StatementEnd
