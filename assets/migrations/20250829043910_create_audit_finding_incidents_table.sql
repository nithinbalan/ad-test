-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_finding_incidents
(
    id INT PRIMARY KEY IDENTITY(1, 1),
	audit_finding_id INT NOT NULL,
    incident_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
	CONSTRAINT FK_audit_finding_incidents_audit_findings FOREIGN KEY (audit_finding_id) REFERENCES audit_findings(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_finding_incidents_incident_register FOREIGN KEY (incident_register_id) REFERENCES incident_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	CONSTRAINT FK_audit_finding_incidents_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
   
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_finding_incidents;

-- +goose StatementEnd