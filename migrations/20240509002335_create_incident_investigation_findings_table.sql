-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_investigation_findings
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    incident_investigation_id INT NOT NULL,
    title NVARCHAR(500) NOT NULL,
    created_by INT NOT NULL,
    updated_by INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_incident_investigation_findings_incident_investigations FOREIGN KEY (incident_investigation_id) REFERENCES incident_investigations(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_investigation_findings_users_1 FOREIGN KEY (created_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_investigation_findings_users_2 FOREIGN KEY (updated_by) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_investigation_findings;

-- +goose StatementEnd