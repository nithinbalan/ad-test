-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_evidences (
    id INT PRIMARY KEY IDENTITY(1, 1),
    incident_register_id INT NOT NULL,
    incident_evidence_source_id INT NOT NULL,
    incident_evidence_type_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) DEFAULT NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) DEFAULT NULL,
    type VARCHAR(255) DEFAULT NULL,
    is_read_only BIT NOT NULL DEFAULT 0,
    is_download BIT NOT NULL DEFAULT 0,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_incident_evidences_incident_register
        FOREIGN KEY (incident_register_id) REFERENCES incident_register(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_evidences_incident_evidence_sources
        FOREIGN KEY (incident_evidence_source_id) REFERENCES incident_evidence_sources(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_evidences_incident_evidence_types
        FOREIGN KEY (incident_evidence_type_id) REFERENCES incident_evidence_types(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_evidences_users_1
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_incident_evidences_users_2
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_evidences;
-- +goose StatementEnd
