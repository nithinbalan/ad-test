-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_requirement_status_update_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_requirement_status_update_id INT NOT NULL,
    url VARCHAR(255) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    type VARCHAR(255) DEFAULT NULL,
    icon VARCHAR(10) NOT NULL, -- Assuming icon names like 'pdf', 'excel', etc.
        CONSTRAINT CHK_compliance_requirement_status_update_documents_icon CHECK (icon IN ('pdf', 'excel', 'ppt', 'image', 'text', 'other')),
    created_by INT,
    updated_by INT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_compliance_requirement_status_update_documents FOREIGN KEY (compliance_requirement_status_update_id) REFERENCES compliance_requirement_status_updates(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirement_status_update_documents_user_1 FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirement_status_update_documents_user_2 FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.compliance_requirement_status_update_documents;
-- +goose StatementEnd
