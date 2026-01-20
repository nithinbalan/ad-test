-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_report_documents (
    id INT PRIMARY KEY IDENTITY(1,1),

    incident_document_id INT NOT NULL,
    incident_report_incident_id INT NOT NULL,
    incident_investigation_id INT NULL,

    url VARCHAR(1000) NULL,
    name NVARCHAR(255) NULL,
    size INT NULL,
    token VARCHAR(255) NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) NULL,

    created_by_id INT NULL,
    updated_by_id INT NULL,
    
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_incident_report_documents_incident_documents FOREIGN KEY (incident_document_id) REFERENCES incident_documents (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_report_documents_incident_report_incidents FOREIGN KEY (incident_report_incident_id) REFERENCES incident_report_incidents (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_report_documents_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_incident_report_documents_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_report_documents;
-- +goose StatementEnd
