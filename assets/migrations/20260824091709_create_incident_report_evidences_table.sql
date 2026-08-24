-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_report_evidences (
    id INT PRIMARY KEY IDENTITY(1, 1),
    incident_evidence_id INT NOT NULL,
    incident_report_incident_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) DEFAULT NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) DEFAULT NULL,
    type VARCHAR(255) DEFAULT NULL,
    created_by_id INT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_incident_report_evidences_incident_evidences
        FOREIGN KEY (incident_evidence_id) REFERENCES incident_evidences(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_report_evidences_incident_report_incidents
        FOREIGN KEY (incident_report_incident_id) REFERENCES incident_report_incidents(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_report_evidences_users_1
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_report_evidences_users_2
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_report_evidences;
-- +goose StatementEnd
