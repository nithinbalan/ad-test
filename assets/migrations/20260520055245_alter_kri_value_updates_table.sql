-- +goose Up
-- +goose StatementBegin
-- Change value column from INT to DECIMAL(18,2)
ALTER TABLE kri_value_updates
ALTER COLUMN value DECIMAL(18,2) NULL;

-- Add new column
ALTER TABLE kri_value_updates
ADD kri_frequency_value_id INT NULL;

-- Add new foreign key constraint
ALTER TABLE kri_value_updates
ADD CONSTRAINT FK_kri_value_updates_kri_frequency_values
FOREIGN KEY (kri_frequency_value_id)
REFERENCES kri_frequency_values(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Drop new foreign key constraint
ALTER TABLE kri_value_updates
DROP CONSTRAINT FK_kri_value_updates_kri_frequency_values;

-- Drop new column
ALTER TABLE kri_value_updates
DROP COLUMN kri_frequency_value_id;

-- Revert value column from DECIMAL(18,2) to INT
ALTER TABLE kri_value_updates
ALTER COLUMN value INT NULL;
-- +goose StatementEnd