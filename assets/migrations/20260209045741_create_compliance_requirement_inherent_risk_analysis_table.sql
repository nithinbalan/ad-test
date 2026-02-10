-- +goose Up
-- +goose StatementBegin
CREATE TABLE compliance_requirement_inherent_risk_analysis 
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    compliance_requirement_id INT NOT NULL,
    compliance_assessment_impact_type_id INT NOT NULL,
    score INT NOT NULL,
    comment NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_compliance_requirement_inherent_risk_analysis_compliance_requirement FOREIGN KEY (compliance_requirement_id) REFERENCES compliance_requirements(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirement_inherent_risk_analysis_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirement_inherent_risk_analysis_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_compliance_requirement_inherent_risk_analysis_impact_types FOREIGN KEY (compliance_assessment_impact_type_id) REFERENCES compliance_assessment_impact_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE compliance_requirement_inherent_risk_analysis;
-- +goose StatementEnd
