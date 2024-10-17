-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_final_report_sections
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    reference_number VARCHAR(45) NOT NULL,
    audit_final_report_id INT NOT NULL,
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    type VARCHAR(20) NOT NULL,
    CONSTRAINT CHK_audit_final_report_sections_type CHECK (type IN ('overview', 'finding-summary', 'findings', 'conclusion', 'custom')),
    [order] DECIMAL(10, 2) NULL,
    audit_final_report_section_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,

    CONSTRAINT FK_audit_final_report_sections_audit_final_reports FOREIGN KEY (audit_final_report_id) REFERENCES audit_final_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_final_report_sections_self FOREIGN KEY (audit_final_report_section_id) REFERENCES audit_final_report_sections(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_final_report_sections_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_final_report_sections_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_final_report_sections;
-- +goose StatementEnd
