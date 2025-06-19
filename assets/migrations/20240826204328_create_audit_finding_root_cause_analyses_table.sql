-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_finding_root_cause_analyses (
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_finding_id INT NOT NULL,
    why INT NOT NULL,
    question NVARCHAR(500) NOT NULL,
    answer NVARCHAR(1000) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_audit_finding_root_cause_analyses FOREIGN KEY (audit_finding_id) REFERENCES audit_findings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_finding_root_cause_analyses_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_finding_root_cause_analyses_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.audit_finding_root_cause_analyses;
-- +goose StatementEnd
