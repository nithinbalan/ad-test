-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_assessment_requirement_responsible_units
ADD residual_risk_score DECIMAL(10,2) NULL;

ALTER TABLE compliance_assessment_requirement_responsible_units
ADD residual_risk_rating_id INT NULL;

ALTER TABLE compliance_assessment_requirement_responsible_units
ADD CONSTRAINT FK_compliance_assessment_requirement_responsible_units_residual_risk_rating FOREIGN KEY (residual_risk_rating_id) REFERENCES compliance_assessment_ratings(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE compliance_assessment_requirement_responsible_units
DROP CONSTRAINT FK_compliance_assessment_requirement_responsible_units_residual_risk_rating;

ALTER TABLE compliance_assessment_requirement_responsible_units
DROP COLUMN residual_risk_rating_id;

ALTER TABLE compliance_assessment_requirement_responsible_units
DROP COLUMN residual_risk_score;

-- +goose StatementEnd
