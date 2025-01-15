-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_final_report_appendices
(
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_final_report_id INT NOT NULL,
    audit_final_report_section_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_audit_final_report_appendices_final_reports FOREIGN KEY (audit_final_report_id) REFERENCES audit_final_reports(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_final_report_appendices_audit_final_report_sections FOREIGN KEY (audit_final_report_section_id) REFERENCES audit_final_report_sections(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_final_report_appendices_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_final_report_appendices_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.audit_final_report_appendices;
-- +goose StatementEnd
