-- +goose Up
-- +goose StatementBegin

CREATE TABLE compliance_assessment_frameworks (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description NVARCHAR(MAX) NULL,
    compliance_assessment_framework_type_id INT NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_compliance_assessment_frameworks_status CHECK (status IN ('active', 'inactive')),
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    
    CONSTRAINT FK_compliance_assessment_frameworks_compliance_assessment_framework_types FOREIGN KEY (compliance_assessment_framework_type_id) 
        REFERENCES compliance_assessment_framework_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_frameworks_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_frameworks_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS compliance_assessment_frameworks;

-- +goose StatementEnd