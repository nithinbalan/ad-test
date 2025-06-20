-- +goose Up
-- +goose StatementBegin
ALTER TABLE kpi_register
ADD review_frequency_id INT NULL;

ALTER TABLE kpi_register
ADD calculation_type_id INT NULL;

ALTER TABLE kpi_register
ADD formula NVARCHAR(255) NULL;

ALTER TABLE kpi_register
ADD unit_id INT NULL;

ALTER TABLE kpi_register
ADD source NVARCHAR(MAX) NULL;

ALTER TABLE kpi_register
ADD actual_value FLOAT NULL;

ALTER TABLE kpi_register
ADD target_value FLOAT NULL;

ALTER TABLE kpi_register
ADD low_threshold INT NULL;

ALTER TABLE kpi_register
ADD medium_threshold INT NULL;

ALTER TABLE kpi_register
ADD weightage FLOAT NULL;

ALTER TABLE kpi_register
ADD start_date DATE NULL;

ALTER TABLE kpi_register
ADD end_date DATE NULL;

ALTER TABLE kpi_register
ADD remarks NVARCHAR(MAX) NULL;

ALTER TABLE kpi_register
ADD CONSTRAINT fk_kpi_register_review_frequency
FOREIGN KEY (review_frequency_id) REFERENCES review_frequencies(id);

ALTER TABLE kpi_register
ADD CONSTRAINT fk_kpi_register_calculation_type
FOREIGN KEY (calculation_type_id) REFERENCES calculation_types(id);

ALTER TABLE kpi_register
ADD CONSTRAINT fk_kpi_register_unit
FOREIGN KEY (unit_id) REFERENCES units(id);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
