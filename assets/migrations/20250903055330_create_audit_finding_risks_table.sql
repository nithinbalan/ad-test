-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_finding_risks (
    id INT PRIMARY KEY IDENTITY(1, 1),
    audit_finding_id INT NOT NULL,
    risk_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_audit_finding_risks_audit_finding FOREIGN KEY (audit_finding_id) REFERENCES audit_findings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_finding_risks_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	CONSTRAINT FK_audit_finding_risks_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_finding_risks;
-- +goose StatementEnd