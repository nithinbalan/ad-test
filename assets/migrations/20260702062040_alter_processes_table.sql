-- +goose Up
-- +goose StatementBegin
-- Change foreign key references from time_units to business_continuity_duration_units

-- Drop old foreign keys
ALTER TABLE processes
DROP CONSTRAINT FK_processes_rto_time_units;

ALTER TABLE processes
DROP CONSTRAINT FK_processes_mtps_time_units;

ALTER TABLE processes
DROP CONSTRAINT FK_processes_rpo_time_units;

-- Add new foreign keys
ALTER TABLE processes
ADD CONSTRAINT FK_processes_rto_bc_duration_units
FOREIGN KEY (rto_unit_id)
REFERENCES business_continuity_duration_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

ALTER TABLE processes
ADD CONSTRAINT FK_processes_mtpd_bc_duration_units
FOREIGN KEY (mtpd_unit_id)
REFERENCES business_continuity_duration_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

ALTER TABLE processes
ADD CONSTRAINT FK_processes_rpo_bc_duration_units
FOREIGN KEY (rpo_unit_id)
REFERENCES business_continuity_duration_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- Add condition columns
ALTER TABLE processes
ADD rto_condition VARCHAR(20) NULL,
    mtpd_condition VARCHAR(20) NULL,
    rpo_condition VARCHAR(20) NULL;

-- Add impact_analysis_level_id column and FK
ALTER TABLE processes
ADD impact_analysis_level_id INT NULL;

ALTER TABLE processes
ADD CONSTRAINT FK_processes_impact_analysis_level
FOREIGN KEY (impact_analysis_level_id)
REFERENCES bcm_impact_analysis_levels(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- Revert foreign keys back to time_units

-- Drop new foreign keys
ALTER TABLE processes
DROP CONSTRAINT FK_processes_rto_bc_duration_units;

ALTER TABLE processes
DROP CONSTRAINT FK_processes_mtpd_bc_duration_units;

ALTER TABLE processes
DROP CONSTRAINT FK_processes_rpo_bc_duration_units;

-- Drop impact_analysis_level_id FK and column
ALTER TABLE processes
DROP CONSTRAINT FK_processes_impact_analysis_level;

ALTER TABLE processes
DROP COLUMN impact_analysis_level_id;

-- Drop condition columns
ALTER TABLE processes
DROP COLUMN rto_condition, mtpd_condition, rpo_condition;

-- Restore old foreign keys
ALTER TABLE processes
ADD CONSTRAINT FK_processes_rto_time_units
FOREIGN KEY (rto_unit_id)
REFERENCES time_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

ALTER TABLE processes
ADD CONSTRAINT FK_processes_mtps_time_units
FOREIGN KEY (mtpd_unit_id)
REFERENCES time_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

ALTER TABLE processes
ADD CONSTRAINT FK_processes_rpo_time_units
FOREIGN KEY (rpo_unit_id)
REFERENCES time_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;
-- +goose StatementEnd