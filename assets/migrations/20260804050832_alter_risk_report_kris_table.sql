-- +goose Up
-- +goose StatementBegin

ALTER TABLE risk_report_kris
ALTER COLUMN current_value decimal(18,2) NULL;

ALTER TABLE risk_report_kris
ADD annual_target decimal(18,2) NULL;

ALTER TABLE risk_report_kris
ADD minimum_susceptibility_percent decimal(18,2) NULL;

ALTER TABLE risk_report_kris
ADD maximum_susceptibility_percent decimal(18,2) NULL;

ALTER TABLE risk_report_kris
ADD maximum_tolerable_percent decimal(18,2) NULL;

ALTER TABLE risk_report_kris
ADD is_lower_better BIT NOT NULL
        CONSTRAINT DF_risk_report_kris_is_lower_better DEFAULT (0);

ALTER TABLE risk_report_kris
ADD is_higher_better BIT NOT NULL
        CONSTRAINT DF_risk_report_kris_is_higher_better DEFAULT (0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_report_kris
DROP CONSTRAINT DF_risk_report_kris_is_lower_better;

ALTER TABLE risk_report_kris
DROP CONSTRAINT DF_risk_report_kris_is_higher_better;

ALTER TABLE risk_report_kris
DROP COLUMN annual_target;

ALTER TABLE risk_report_kris
DROP COLUMN minimum_susceptibility_percent;

ALTER TABLE risk_report_kris
DROP COLUMN maximum_susceptibility_percent;

ALTER TABLE risk_report_kris
DROP COLUMN maximum_tolerable_percent;

ALTER TABLE risk_report_kris
DROP COLUMN is_lower_better;

ALTER TABLE risk_report_kris
DROP COLUMN is_higher_better;

ALTER TABLE risk_report_kris
ALTER COLUMN current_value INT NULL;
-- +goose StatementEnd
