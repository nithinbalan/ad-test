-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_investigation_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    incident_investigation_id INT NOT NULL,
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
    CONSTRAINT FK_incident_investigation_documents_incident_investigations FOREIGN KEY (incident_investigation_id) REFERENCES incident_investigations(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_investigation_documents_users_1 FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_investigation_documents_users_2 FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.incident_investigation_documents;
-- +goose StatementEnd
