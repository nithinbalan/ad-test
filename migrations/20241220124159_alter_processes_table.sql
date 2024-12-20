-- +goose Up
-- +goose StatementBegin
-- Add the recovery objective columns with a foreign key to the time_units table

ALTER TABLE processes 
ADD rto INT NULL;

ALTER TABLE processes 
ADD rto_unit_id INT NULL;

ALTER TABLE processes 
ADD mtpd INT NULL;

ALTER TABLE processes 
ADD mtpd_unit_id INT NULL;

ALTER TABLE processes 
ADD rpo INT NULL;

ALTER TABLE processes 
ADD rpo_unit_id INT NULL;

-- Add the foreign key constraints for rto_unit_id, mtpd_unit_id, rpo_unit_id
ALTER TABLE processes
ADD CONSTRAINT FK_processes_rto_time_units FOREIGN KEY (rto_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE processes
ADD CONSTRAINT FK_processes_mtps_time_units FOREIGN KEY (mtpd_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE processes
ADD CONSTRAINT FK_processes_rpo_time_units FOREIGN KEY (rpo_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop foreign key constraints before removing columns
ALTER TABLE processes
DROP CONSTRAINT FK_processes_rto_time_units;

ALTER TABLE processes
DROP COLUMN rto_unit_id;

ALTER TABLE processes
DROP CONSTRAINT FK_processes_mtps_time_units;

ALTER TABLE processes
DROP COLUMN mtpd_unit_id;

ALTER TABLE processes
DROP CONSTRAINT FK_processes_rpo_time_units;

ALTER TABLE processes
DROP COLUMN rpo_unit_id;

-- Remove the additional columns
ALTER TABLE processes 
DROP COLUMN rto;

ALTER TABLE processes 
DROP COLUMN mtpd;

ALTER TABLE processes 
DROP COLUMN rpo;
-- +goose StatementEnd
