-- +goose Up
-- +goose StatementBegin

ALTER TABLE compliance_report_requirements
DROP CONSTRAINT fk_compliance_report_requirements_inherent_risk_rating;

ALTER TABLE compliance_report_requirements
DROP CONSTRAINT fk_compliance_report_requirements_residual_risk_rating;

ALTER TABLE compliance_report_requirements
ADD CONSTRAINT fk_compliance_report_requirements_inherent_risk_rating
    FOREIGN KEY (inherent_risk_rating_id)
    REFERENCES compliance_assessment_ratings(id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION;

ALTER TABLE compliance_report_requirements
ADD CONSTRAINT fk_compliance_report_requirements_residual_risk_rating
    FOREIGN KEY (residual_risk_rating_id)
    REFERENCES compliance_assessment_ratings(id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE compliance_report_requirements
DROP CONSTRAINT fk_compliance_report_requirements_inherent_risk_rating;

ALTER TABLE compliance_report_requirements
DROP CONSTRAINT fk_compliance_report_requirements_residual_risk_rating;

ALTER TABLE compliance_report_requirements
ADD CONSTRAINT fk_compliance_report_requirements_inherent_risk_rating
    FOREIGN KEY (inherent_risk_rating_id)
    REFERENCES compliance_assessment_ratings(id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION;

ALTER TABLE compliance_report_requirements
ADD CONSTRAINT fk_compliance_report_requirements_residual_risk_rating
    FOREIGN KEY (residual_risk_rating_id)
    REFERENCES compliance_assessment_ratings(id)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION;

-- +goose StatementEnd