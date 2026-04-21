-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_register
ADD average_residual_risk_score DECIMAL(10,2) NULL;

ALTER TABLE compliance_register
ADD average_residual_risk_rating_id INT NULL;

ALTER TABLE compliance_register
ADD CONSTRAINT FK_compliance_register_residual_risk_rating FOREIGN KEY (average_residual_risk_rating_id) REFERENCES compliance_assessment_ratings(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE compliance_register
DROP CONSTRAINT FK_compliance_register_residual_risk_rating;

ALTER TABLE compliance_register
DROP COLUMN average_residual_risk_rating_id;

ALTER TABLE compliance_register
DROP COLUMN average_residual_risk_score;

-- +goose StatementEnd
