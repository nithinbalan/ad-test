-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_assessment_requirements (
    id INT PRIMARY KEY IDENTITY(1,1),
    compliance_assessment_id INT NOT NULL,
    compliance_requirement_id INT NOT NULL,
    compliance_status_id INT NULL,
    score INT NULL,
    comments NVARCHAR(MAX) NULL,
    is_completed TINYINT NOT NULL DEFAULT 0,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_compliance_assessment_requirements_compliance_assessments FOREIGN KEY (compliance_assessment_id) REFERENCES compliance_assessments(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_requirements_compliance_requirements FOREIGN KEY (compliance_requirement_id) REFERENCES compliance_requirements(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_requirements_compliance_statuses FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_requirements_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_assessment_requirements_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.compliance_assessment_requirements;
-- +goose StatementEnd
