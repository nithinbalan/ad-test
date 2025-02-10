-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_assessment_requirement_required_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_assessment_requirement_document_title_id INT NOT NULL,
    language VARCHAR(10) NOT NULL, -- Assuming language codes like 'en' and 'ar'
        CONSTRAINT CHK_compliance_assessment_requirement_required_documents_language CHECK (language IN ('en', 'ar')),
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    type VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    
    CONSTRAINT FK_compliance_assessment_requirement_required_documents_compliance_assessment_requirement_document_titles FOREIGN KEY (compliance_assessment_requirement_document_title_id) REFERENCES compliance_assessment_requirement_document_titles(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_requirement_required_documents_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_requirement_required_documents_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.compliance_assessment_requirement_required_documents;
-- +goose StatementEnd