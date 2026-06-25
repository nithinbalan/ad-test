-- +goose Up
-- +goose StatementBegin

-- Alter processes table
-- Alter columns to nullable
ALTER TABLE processes ALTER COLUMN process_group_id INT NULL;
ALTER TABLE processes ALTER COLUMN lifecycle_status_id INT NULL;

-- Add new columns to processes table
ALTER TABLE processes ADD division_id INT NULL;
ALTER TABLE processes ADD is_automatic BIT NOT NULL CONSTRAINT DF_processes_is_automatic DEFAULT 0;
ALTER TABLE processes ADD is_manual BIT NOT NULL CONSTRAINT DF_processes_is_manual DEFAULT 0;
ALTER TABLE processes ADD process_type_id INT NULL;
ALTER TABLE processes ADD location_id INT NULL;
ALTER TABLE processes ADD alternative_location NVARCHAR(MAX) NULL;
ALTER TABLE processes ADD rush_hour VARCHAR(100) NULL;
ALTER TABLE processes ADD review_frequency_id INT NULL;
ALTER TABLE processes ADD is_can_be_done_manually BIT NOT NULL CONSTRAINT DF_processes_is_can_be_done_manually DEFAULT 0;
ALTER TABLE processes ADD time_to_perform_procedure_manually VARCHAR(100) NULL;
ALTER TABLE processes ADD steps_to_perform_procedure_manually NVARCHAR(MAX) NULL;
ALTER TABLE processes ADD requirement_to_perform_procedure_manually NVARCHAR(MAX) NULL;


-- Add foreign key constraints for new columns
ALTER TABLE processes ADD CONSTRAINT FK_processes_division FOREIGN KEY (division_id) REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE processes ADD CONSTRAINT FK_processes_bcm_process_types FOREIGN KEY (process_type_id) REFERENCES bcm_process_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE processes ADD CONSTRAINT FK_processes_locations_1 FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
ALTER TABLE processes ADD CONSTRAINT FK_processes_review_frequencies FOREIGN KEY (review_frequency_id) REFERENCES review_frequencies(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop new columns from processes table
ALTER TABLE processes DROP CONSTRAINT FK_processes_division;
ALTER TABLE processes DROP CONSTRAINT FK_processes_review_frequencies;
ALTER TABLE processes DROP CONSTRAINT FK_processes_locations_1;
ALTER TABLE processes DROP CONSTRAINT FK_processes_bcm_process_types;
ALTER TABLE processes DROP COLUMN review_frequency_id;
ALTER TABLE processes DROP COLUMN rush_hour;
ALTER TABLE processes DROP COLUMN alternative_location;
ALTER TABLE processes DROP COLUMN location_id;
ALTER TABLE processes DROP COLUMN process_type_id;
ALTER TABLE processes DROP CONSTRAINT DF_processes_is_manual;
ALTER TABLE processes DROP COLUMN is_manual;
ALTER TABLE processes DROP CONSTRAINT DF_processes_is_automatic;
ALTER TABLE processes DROP CONSTRAINT DF_processes_is_can_be_done_manually;
ALTER TABLE processes DROP COLUMN is_automatic;
ALTER TABLE processes DROP COLUMN division_id;
ALTER TABLE processes DROP COLUMN requirement_to_perform_procedure_manually;
ALTER TABLE processes DROP COLUMN steps_to_perform_procedure_manually;
ALTER TABLE processes DROP COLUMN time_to_perform_procedure_manually;
ALTER TABLE processes DROP COLUMN is_can_be_done_manually;
-- +goose StatementEnd
