-- +goose Up
-- +goose StatementBegin
-- Drop old foreign keys referencing time_units
ALTER TABLE bcp_processes
DROP CONSTRAINT FK_bcp_processes_time_units_1;

ALTER TABLE bcp_processes
DROP CONSTRAINT FK_bcp_processes_time_units_2;

ALTER TABLE bcp_processes
DROP CONSTRAINT FK_bcp_processes_time_units_3;

-- Add new foreign keys referencing business_continuity_duration_units
ALTER TABLE bcp_processes
ADD CONSTRAINT FK_bcp_processes_rto_bc_duration_units
FOREIGN KEY (rto_unit_id)
REFERENCES business_continuity_duration_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

ALTER TABLE bcp_processes
ADD CONSTRAINT FK_bcp_processes_mtpd_bc_duration_units
FOREIGN KEY (mtpd_unit_id)
REFERENCES business_continuity_duration_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

ALTER TABLE bcp_processes
ADD CONSTRAINT FK_bcp_processes_rpo_bc_duration_units
FOREIGN KEY (rpo_unit_id)
REFERENCES business_continuity_duration_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- Add condition columns
ALTER TABLE bcp_processes
ADD rto_condition VARCHAR(20) NULL,
    mtpd_condition VARCHAR(20) NULL,
    rpo_condition VARCHAR(20) NULL;

-- Add impact_analysis_level_id column and FK
ALTER TABLE bcp_processes
ADD impact_analysis_level_id INT NULL;

ALTER TABLE bcp_processes
ADD CONSTRAINT FK_bcp_processes_impact_analysis_level
FOREIGN KEY (impact_analysis_level_id)
REFERENCES bcm_impact_analysis_levels(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- Drop new foreign keys
ALTER TABLE bcp_processes
DROP CONSTRAINT FK_bcp_processes_rto_bc_duration_units;

ALTER TABLE bcp_processes
DROP CONSTRAINT FK_bcp_processes_mtpd_bc_duration_units;

ALTER TABLE bcp_processes
DROP CONSTRAINT FK_bcp_processes_rpo_bc_duration_units;

-- Drop impact_analysis_level_id FK and column
ALTER TABLE bcp_processes
DROP CONSTRAINT FK_bcp_processes_impact_analysis_level;

ALTER TABLE bcp_processes
DROP COLUMN impact_analysis_level_id;

-- Drop condition columns
ALTER TABLE bcp_processes
DROP COLUMN rto_condition, mtpd_condition, rpo_condition;

-- Restore old foreign keys referencing time_units
ALTER TABLE bcp_processes
ADD CONSTRAINT FK_bcp_processes_time_units_1
FOREIGN KEY (rto_unit_id)
REFERENCES time_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

ALTER TABLE bcp_processes
ADD CONSTRAINT FK_bcp_processes_time_units_2
FOREIGN KEY (mtpd_unit_id)
REFERENCES time_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

ALTER TABLE bcp_processes
ADD CONSTRAINT FK_bcp_processes_time_units_3
FOREIGN KEY (rpo_unit_id)
REFERENCES time_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;
-- +goose StatementEnd
