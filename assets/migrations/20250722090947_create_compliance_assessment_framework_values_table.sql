-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_assessment_framework_values (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_assessment_framework_id INT NOT NULL,
    compliance_status_id INT NOT NULL,
    value INT NULL,
    title_en NVARCHAR(255) NULL,
    title_ar NVARCHAR(255) NULL,
    description NVARCHAR(MAX) NULL,
    score_min INT NULL,
    score_max INT NULL,
    is_applicable TINYINT NOT NULL DEFAULT 1,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_compliance_assessment_framework_values_compliance_assessment_frameworks FOREIGN KEY (compliance_assessment_framework_id)
        REFERENCES compliance_assessment_frameworks(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_framework_values_compliance_statuses FOREIGN KEY (compliance_status_id)
        REFERENCES compliance_statuses(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_framework_values_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_framework_values_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS compliance_assessment_framework_values;
-- +goose StatementEnd 