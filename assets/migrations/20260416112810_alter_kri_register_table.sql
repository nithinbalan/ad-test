-- +goose Up
-- +goose StatementBegin

ALTER TABLE kri_register
ADD 
    level_of_measurement_id INT NULL,
    risk_responsible_unit_id INT NULL,
    risk_sub_category_id INT NULL,
    annual_target DECIMAL(18,2) NULL,
    review_target DECIMAL(18,2) NULL,
    distribution_status_id INT NULL,
    minimum_susceptibility_percent DECIMAL(5,2) NULL,
    maximum_susceptibility_percent DECIMAL(5,2) NULL,
    maximum_tolerable_percent DECIMAL(5,2) NULL;

-- Foreign Keys
ALTER TABLE kri_register
ADD CONSTRAINT FK_kri_register_level_of_measurements
FOREIGN KEY (level_of_measurement_id)
REFERENCES level_of_measurements(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE kri_register
ADD CONSTRAINT FK_kri_register_risk_responsible_unit
FOREIGN KEY (risk_responsible_unit_id)
REFERENCES organization_units(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE kri_register
ADD CONSTRAINT FK_kri_register_risk_sub_categories
FOREIGN KEY (risk_sub_category_id)
REFERENCES risk_sub_categories(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE kri_register
ADD CONSTRAINT FK_kri_register_distribution_statuses
FOREIGN KEY (distribution_status_id)
REFERENCES distribution_statuses(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop constraints first
ALTER TABLE kri_register DROP CONSTRAINT FK_kri_register_distribution_statuses;
ALTER TABLE kri_register DROP CONSTRAINT FK_kri_register_risk_sub_categories;
ALTER TABLE kri_register DROP CONSTRAINT FK_kri_register_risk_responsible_unit;
ALTER TABLE kri_register DROP CONSTRAINT FK_kri_register_level_of_measurements;

-- Drop columns
ALTER TABLE kri_register
DROP COLUMN 
    level_of_measurement_id,
    risk_responsible_unit_id,
    risk_sub_category_id,
    annual_target,
    review_target,
    distribution_status_id,
    minimum_susceptibility_percent,
    maximum_susceptibility_percent,
    maximum_tolerable_percent;

-- +goose StatementEnd