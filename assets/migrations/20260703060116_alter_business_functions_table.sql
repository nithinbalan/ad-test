-- +goose Up
-- +goose StatementBegin
-- Drop old foreign keys referencing time_units
ALTER TABLE business_functions
DROP CONSTRAINT FK_business_functions_rto_time_units;

ALTER TABLE business_functions
DROP CONSTRAINT FK_business_functions_mtpd_time_units;

ALTER TABLE business_functions
DROP CONSTRAINT FK_business_functions_rpo_time_units;

-- Add new foreign keys referencing business_continuity_duration_units
ALTER TABLE business_functions
ADD CONSTRAINT FK_business_functions_rto_bc_duration_units
FOREIGN KEY (rto_unit_id)
REFERENCES business_continuity_duration_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

ALTER TABLE business_functions
ADD CONSTRAINT FK_business_functions_mtpd_bc_duration_units
FOREIGN KEY (mtpd_unit_id)
REFERENCES business_continuity_duration_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

ALTER TABLE business_functions
ADD CONSTRAINT FK_business_functions_rpo_bc_duration_units
FOREIGN KEY (rpo_unit_id)
REFERENCES business_continuity_duration_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- Add condition columns
ALTER TABLE business_functions
ADD rto_condition VARCHAR(20) NULL,
    mtpd_condition VARCHAR(20) NULL,
    rpo_condition VARCHAR(20) NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- Drop new foreign keys
ALTER TABLE business_functions
DROP CONSTRAINT FK_business_functions_rto_bc_duration_units;

ALTER TABLE business_functions
DROP CONSTRAINT FK_business_functions_mtpd_bc_duration_units;

ALTER TABLE business_functions
DROP CONSTRAINT FK_business_functions_rpo_bc_duration_units;

-- Drop condition columns
ALTER TABLE business_functions
DROP COLUMN rto_condition, mtpd_condition, rpo_condition;

-- Restore old foreign keys referencing time_units
ALTER TABLE business_functions
ADD CONSTRAINT FK_business_functions_rto_time_units
FOREIGN KEY (rto_unit_id)
REFERENCES time_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

ALTER TABLE business_functions
ADD CONSTRAINT FK_business_functions_mtpd_time_units
FOREIGN KEY (mtpd_unit_id)
REFERENCES time_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

ALTER TABLE business_functions
ADD CONSTRAINT FK_business_functions_rpo_time_units
FOREIGN KEY (rpo_unit_id)
REFERENCES time_units(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;
-- +goose StatementEnd
