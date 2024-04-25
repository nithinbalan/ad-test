-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_register_id INT NOT NULL,
    language VARCHAR(10) NOT NULL, -- Assuming language codes like 'en' and 'ar'
        CONSTRAINT CHK_compliance_documents_language CHECK (language IN ('en', 'ar')),
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    type VARCHAR(255) DEFAULT NULL,
    icon VARCHAR(10) NOT NULL, -- Assuming icon names like 'pdf', 'excel', etc.
        CONSTRAINT CHK_compliance_documents_icon CHECK (icon IN ('pdf', 'excel', 'ppt', 'image', 'text', 'other')),
    created_by INT,
    updated_by INT NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT FK_compliance_documents_compliance_register FOREIGN KEY (compliance_register_id) REFERENCES compliance_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_documents_users_1 FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_documents_users_2 FOREIGN KEY (updated_by) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.compliance_documents;
-- +goose StatementEnd
