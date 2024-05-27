-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    incident_register_id INT NOT NULL,
    language VARCHAR(10) NOT NULL,
        CONSTRAINT CHK_incident_documents_language CHECK (language IN ('en', 'ar')),
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,
    icon VARCHAR(10) NOT NULL,
    created_by INT,
    updated_by INT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_incident_documents_incident_register FOREIGN KEY (incident_register_id) REFERENCES incident_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_documents_users_1 FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_documents_users_2 FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.incident_documents;
-- +goose StatementEnd
