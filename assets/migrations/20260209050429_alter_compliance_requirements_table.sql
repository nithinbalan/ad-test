-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_requirements
ADD inherent_risk_score DECIMAL(10,2) NULL;

ALTER TABLE compliance_requirements
ADD residual_risk_score DECIMAL(10,2) NULL;

ALTER TABLE compliance_requirements
ADD inherent_risk_rating_id INT NULL;

ALTER TABLE compliance_requirements
ADD residual_risk_rating_id INT NULL;

ALTER TABLE compliance_requirements
ADD CONSTRAINT FK_compliance_requirements_inherent_risk_rating FOREIGN KEY (inherent_risk_rating_id) REFERENCES compliance_assessment_ratings(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE compliance_requirements
ADD CONSTRAINT FK_compliance_requirements_residual_risk_rating FOREIGN KEY (residual_risk_rating_id) REFERENCES compliance_assessment_ratings(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE compliance_requirements
DROP CONSTRAINT FK_compliance_requirements_inherent_risk_rating;

ALTER TABLE compliance_requirements
DROP CONSTRAINT FK_compliance_requirements_residual_risk_rating;

ALTER TABLE compliance_requirements
DROP COLUMN inherent_risk_rating_id;

ALTER TABLE compliance_requirements
DROP COLUMN inherent_risk_score;

ALTER TABLE compliance_requirements
DROP COLUMN residual_risk_rating_id;

ALTER TABLE compliance_requirements
DROP COLUMN residual_risk_score;

-- +goose StatementEnd
