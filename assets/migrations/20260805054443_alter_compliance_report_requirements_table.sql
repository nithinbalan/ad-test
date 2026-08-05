-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_report_requirements
ADD compliance_register_id INT NULL;

ALTER TABLE compliance_report_requirements
ADD CONSTRAINT FK_compliance_report_requirements_compliance_register FOREIGN KEY (compliance_register_id)
REFERENCES compliance_register (id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE compliance_report_requirements
ALTER COLUMN title_en NVARCHAR(MAX) NOT NULL;

ALTER TABLE compliance_report_requirements
ALTER COLUMN title_ar NVARCHAR(MAX) NOT NULL;

ALTER TABLE compliance_report_requirements
ADD inherent_score DECIMAL(10,2) NULL;

ALTER TABLE compliance_report_requirements
ADD residual_score DECIMAL(10,2) NULL;

ALTER TABLE compliance_report_requirements
ADD is_acknowledged BIT NOT NULL
        CONSTRAINT DF_compliance_report_requirements_is_acknowledged DEFAULT (0);

ALTER TABLE compliance_report_requirements
ADD acknowledge_user_id INT NULL;

ALTER TABLE compliance_report_requirements
ADD acknowledge_comment NVARCHAR(MAX) NULL;

ALTER TABLE compliance_report_requirements
ADD is_risk_based_assessment BIT NULL;

ALTER TABLE compliance_report_requirements
ADD is_regulation_notification_assessment BIT NULL;

ALTER TABLE compliance_report_requirements
ADD CONSTRAINT FK_compliance_report_requirements_users_acknowledge FOREIGN KEY (acknowledge_user_id)
REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE compliance_report_requirements
ADD inherent_risk_rating_id INT NULL;

ALTER TABLE compliance_report_requirements
ADD CONSTRAINT fk_compliance_report_requirements_inherent_risk_rating
    FOREIGN KEY (inherent_risk_rating_id)
    REFERENCES risk_matrix_risk_ratings(id)
    ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE compliance_report_requirements
ADD residual_risk_rating_id INT NULL;

ALTER TABLE compliance_report_requirements
ADD CONSTRAINT fk_compliance_report_requirements_residual_risk_rating
    FOREIGN KEY (residual_risk_rating_id)
    REFERENCES risk_matrix_risk_ratings(id)
    ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE compliance_report_requirements
DROP CONSTRAINT IF EXISTS DF_compliance_report_requirements_is_acknowledged;

ALTER TABLE compliance_report_requirements
DROP CONSTRAINT IF EXISTS fk_compliance_report_requirements_residual_risk_rating;

ALTER TABLE compliance_report_requirements
DROP CONSTRAINT IF EXISTS fk_compliance_report_requirements_inherent_risk_rating;

ALTER TABLE compliance_report_requirements
DROP COLUMN IF EXISTS residual_risk_rating_id;

ALTER TABLE compliance_report_requirements
DROP COLUMN IF EXISTS inherent_risk_rating_id;

ALTER TABLE compliance_report_requirements
DROP CONSTRAINT FK_compliance_report_requirements_users_acknowledge;

ALTER TABLE compliance_report_requirements
DROP COLUMN is_risk_based_assessment;

ALTER TABLE compliance_report_requirements
DROP COLUMN is_regulation_notification_assessment;

ALTER TABLE compliance_report_requirements
DROP COLUMN acknowledge_comment;

ALTER TABLE compliance_report_requirements
DROP COLUMN acknowledge_user_id;

ALTER TABLE compliance_report_requirements
DROP COLUMN is_acknowledged;

ALTER TABLE compliance_report_requirements
DROP COLUMN residual_score;

ALTER TABLE compliance_report_requirements
DROP COLUMN inherent_score;

ALTER TABLE compliance_report_requirements
ALTER COLUMN title_en NVARCHAR(500) NOT NULL;

ALTER TABLE compliance_report_requirements
ALTER COLUMN title_ar NVARCHAR(500) NOT NULL;

ALTER TABLE compliance_report_requirements
DROP CONSTRAINT FK_compliance_report_requirements_compliance_register;

ALTER TABLE compliance_report_requirements
DROP COLUMN compliance_register_id;
-- +goose StatementEnd