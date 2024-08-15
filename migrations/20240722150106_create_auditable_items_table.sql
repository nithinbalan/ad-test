-- +goose Up
-- +goose StatementBegin
CREATE TABLE auditable_items
(
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(500) NULL,
    title_ar NVARCHAR(500) NULL,
    process_id INT NULL,
    control_id INT NULL,
    document_id INT NULL,
    incident_register_id INT NULL,
    compliance_requirement_id INT NULL,
    risk_register_id INT NULL,
    user_id INT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    auditor_risk_rating_id INT NOT NULL,
    audit_status_id INT NOT NULL,
    responsible_unit_id INT NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_auditable_items_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_auditable_items_controls FOREIGN KEY (control_id) REFERENCES controls(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_auditable_items_documents FOREIGN KEY (document_id) REFERENCES documents(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_auditable_items_incident_registers FOREIGN KEY (incident_register_id) REFERENCES incident_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_auditable_items_compliance_requirements FOREIGN KEY (compliance_requirement_id) REFERENCES compliance_requirements(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_auditable_items_risk_registers FOREIGN KEY (risk_register_id) REFERENCES risk_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_auditable_items_users_1 FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_auditable_items_auditor_risk_ratings FOREIGN KEY (auditor_risk_rating_id) REFERENCES auditor_risk_ratings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_auditable_items_auditor_statuses FOREIGN KEY (audit_status_id) REFERENCES audit_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_auditable_items_responsible_units FOREIGN KEY (responsible_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_auditable_items_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_auditable_items_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.auditable_items;
-- +goose StatementEnd
