-- +goose Up
-- +goose StatementBegin
ALTER TABLE business_functions
ADD highest_process_criticality_level_id INT NULL;

ALTER TABLE business_functions
ADD CONSTRAINT FK_business_functions_process_criticality_levels
FOREIGN KEY (highest_process_criticality_level_id) REFERENCES process_criticality_levels(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE business_functions
DROP CONSTRAINT FK_business_functions_process_criticality_levels;

ALTER TABLE business_functions
DROP COLUMN highest_process_criticality_level_id;
-- +goose StatementEnd
