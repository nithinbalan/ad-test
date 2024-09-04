-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_draft_reports
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    audit_id INT NOT NULL,
    report_date DATE NOT NULL,
    excecutive_summary NVARCHAR(MAX) NULL,
    conclusion NVARCHAR(MAX) NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,

    CONSTRAINT FK_audit_draft_reports_audits FOREIGN KEY (audit_id) REFERENCES audits(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_draft_reports_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_draft_reports_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_draft_reports;
-- +goose StatementEnd
