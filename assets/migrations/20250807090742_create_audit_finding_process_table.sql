-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_finding_processes (
    id INT PRIMARY KEY IDENTITY(1, 1),
    process_id INT NOT NULL,
    audit_finding_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_audit_finding_processes_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_finding_processes_audit_findings FOREIGN KEY (audit_finding_id) REFERENCES audit_findings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	CONSTRAINT FK_audit_finding_processes_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_finding_processes;

-- +goose StatementEnd
