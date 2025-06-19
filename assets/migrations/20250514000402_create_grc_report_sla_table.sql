-- +goose Up
-- +goose StatementBegin
CREATE TABLE grc_report_sla (
    id INT PRIMARY KEY IDENTITY(1,1),
    grc_report_id INT NOT NULL,
    parent_organization_unit_id INT NOT NULL,
    organization_unit_id INT NOT NULL,
    reference_number NVARCHAR(255) NULL,
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NULL,
    compliance_status_id INT NULL,
    lifecycle_status_id INT NOT NULL,
    issue_date DATE NULL,
    expiry_date DATE NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_grc_report_sla_grc_report FOREIGN KEY (grc_report_id) REFERENCES grc_reports(id)  ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_sla_parent_organization_unit FOREIGN KEY (parent_organization_unit_id) REFERENCES organization_units(id)  ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_sla_organization_unit FOREIGN KEY (organization_unit_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_sla_compliance_status FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_sla_user_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_sla_user_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_sla_lifecycle_status FOREIGN KEY (lifecycle_status_id) REFERENCES lifecycle_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE grc_report_sla;
-- +goose StatementEnd
