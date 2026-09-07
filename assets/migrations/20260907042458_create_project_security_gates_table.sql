-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gates (
    id INT PRIMARY KEY IDENTITY(1, 1),
    project_title_en NVARCHAR(100) NOT NULL,
    project_title_ar NVARCHAR(100) NOT NULL,
    stage_id INT NOT NULL,
    template_version_id INT NOT NULL,
    business_owner_id INT NOT NULL,
    readiness_status_id INT NOT NULL,
    notes NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_project_security_gates_readiness_status
        FOREIGN KEY (readiness_status_id)
        REFERENCES project_security_gate_readiness_statuses(id),

    CONSTRAINT FK_project_security_gates_stage
        FOREIGN KEY (stage_id)
        REFERENCES project_security_gate_stages(id),

    CONSTRAINT FK_project_security_gates_template_version
        FOREIGN KEY (template_version_id)
        REFERENCES project_security_gate_checklist_template_versions(id),

    CONSTRAINT FK_project_security_gates_business_owner
        FOREIGN KEY (business_owner_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_security_gates_created_by
        FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_security_gates_updated_by
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gates;
-- +goose StatementEnd