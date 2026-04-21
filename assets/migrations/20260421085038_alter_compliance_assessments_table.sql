-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_assessments
ADD average_residual_risk_score DECIMAL(10,2) NULL;

ALTER TABLE compliance_assessments
ADD average_residual_risk_rating_id INT NULL;

ALTER TABLE compliance_assessments
ADD CONSTRAINT FK_compliance_assessments_residual_risk_rating FOREIGN KEY (average_residual_risk_rating_id) REFERENCES compliance_assessment_ratings(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE compliance_assessments
DROP CONSTRAINT FK_compliance_assessments_residual_risk_rating;

ALTER TABLE compliance_assessments
DROP COLUMN average_residual_risk_rating_id;

ALTER TABLE compliance_assessments
DROP COLUMN average_residual_risk_score;

-- +goose StatementEnd
