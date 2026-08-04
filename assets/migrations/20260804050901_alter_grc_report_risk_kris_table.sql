-- +goose Up
-- +goose StatementBegin
ALTER TABLE grc_report_risk_kris
ALTER COLUMN current_value decimal(18,2) NULL;

ALTER TABLE grc_report_risk_kris
ADD annual_target decimal(18,2) NULL;

ALTER TABLE grc_report_risk_kris
ADD minimum_susceptibility_percent decimal(18,2) NULL;

ALTER TABLE grc_report_risk_kris
ADD maximum_susceptibility_percent decimal(18,2) NULL;

ALTER TABLE grc_report_risk_kris
ADD maximum_tolerable_percent decimal(18,2) NULL;

ALTER TABLE grc_report_risk_kris
ADD is_lower_better BIT NOT NULL
        CONSTRAINT DF_grc_report_risk_kris_is_lower_better DEFAULT (0);

ALTER TABLE grc_report_risk_kris
ADD is_higher_better BIT NOT NULL
        CONSTRAINT DF_grc_report_risk_kris_is_higher_better DEFAULT (0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE grc_report_risk_kris
DROP CONSTRAINT DF_grc_report_risk_kris_is_lower_better;

ALTER TABLE grc_report_risk_kris
DROP CONSTRAINT DF_grc_report_risk_kris_is_higher_better;

ALTER TABLE grc_report_risk_kris
DROP COLUMN annual_target;

ALTER TABLE grc_report_risk_kris
DROP COLUMN minimum_susceptibility_percent;

ALTER TABLE grc_report_risk_kris
DROP COLUMN maximum_susceptibility_percent;

ALTER TABLE grc_report_risk_kris
DROP COLUMN maximum_tolerable_percent;

ALTER TABLE grc_report_risk_kris
DROP COLUMN is_lower_better;

ALTER TABLE grc_report_risk_kris
DROP COLUMN is_higher_better;

ALTER TABLE grc_report_risk_kris
ALTER COLUMN current_value INT NULL;
-- +goose StatementEnd
