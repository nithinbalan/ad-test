-- +goose Up
-- +goose StatementBegin
-- Add the recovery objective columns with a foreign key to the time_units table

ALTER TABLE process_dependencies 
ADD rto INT NULL;

ALTER TABLE process_dependencies 
ADD rto_unit_id INT NULL;

ALTER TABLE process_dependencies 
ADD mtpd INT NULL;

ALTER TABLE process_dependencies 
ADD mtpd_unit_id INT NULL;

ALTER TABLE process_dependencies 
ADD rpo INT NULL;

ALTER TABLE process_dependencies 
ADD rpo_unit_id INT NULL;

-- Add the foreign key constraints for rto_unit_id, mtpd_unit_id, rpo_unit_id
ALTER TABLE process_dependencies
ADD CONSTRAINT FK_process_dependencies_rto_time_units FOREIGN KEY (rto_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE process_dependencies
ADD CONSTRAINT FK_process_dependencies_mtps_time_units FOREIGN KEY (mtpd_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE process_dependencies
ADD CONSTRAINT FK_process_dependencies_rpo_time_units FOREIGN KEY (rpo_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop foreign key constraints before removing columns
ALTER TABLE process_dependencies
DROP CONSTRAINT FK_process_dependencies_rto_time_units;

ALTER TABLE process_dependencies
DROP COLUMN rto_unit_id;

ALTER TABLE process_dependencies
DROP CONSTRAINT FK_process_dependencies_mtps_time_units;

ALTER TABLE process_dependencies
DROP COLUMN mtpd_unit_id;

ALTER TABLE process_dependencies
DROP CONSTRAINT FK_process_dependencies_rpo_time_units;

ALTER TABLE process_dependencies
DROP COLUMN rpo_unit_id;

-- Remove the additional columns
ALTER TABLE process_dependencies 
DROP COLUMN rto;

ALTER TABLE process_dependencies 
DROP COLUMN mtpd;

ALTER TABLE process_dependencies 
DROP COLUMN rpo;
-- +goose StatementEnd
