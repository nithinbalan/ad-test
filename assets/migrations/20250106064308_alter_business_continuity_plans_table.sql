-- +goose Up
-- +goose StatementBegin
-- Add the recovery objective columns to the business_continuity_plans table
ALTER TABLE business_continuity_plans 
ADD rto INT NULL;

ALTER TABLE business_continuity_plans 
ADD rto_unit_id INT NULL;

ALTER TABLE business_continuity_plans 
ADD mtpd INT NULL;

ALTER TABLE business_continuity_plans 
ADD mtpd_unit_id INT NULL;

ALTER TABLE business_continuity_plans 
ADD rpo INT NULL;

ALTER TABLE business_continuity_plans 
ADD rpo_unit_id INT NULL;

-- Add the foreign key constraints for rto_unit_id, mtpd_unit_id, and rpo_unit_id
ALTER TABLE business_continuity_plans
ADD CONSTRAINT FK_business_continuity_plans_rto_time_units FOREIGN KEY (rto_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE business_continuity_plans
ADD CONSTRAINT FK_business_continuity_plans_mtpd_time_units FOREIGN KEY (mtpd_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE business_continuity_plans
ADD CONSTRAINT FK_business_continuity_plans_rpo_time_units FOREIGN KEY (rpo_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove foreign key constraints before dropping columns
ALTER TABLE business_continuity_plans
DROP CONSTRAINT FK_business_continuity_plans_rto_time_units;

ALTER TABLE business_continuity_plans
DROP COLUMN rto_unit_id;

ALTER TABLE business_continuity_plans
DROP CONSTRAINT FK_business_continuity_plans_mtpd_time_units;

ALTER TABLE business_continuity_plans
DROP COLUMN mtpd_unit_id;

ALTER TABLE business_continuity_plans
DROP CONSTRAINT FK_business_continuity_plans_rpo_time_units;

ALTER TABLE business_continuity_plans
DROP COLUMN rpo_unit_id;

-- Remove the additional columns
ALTER TABLE business_continuity_plans 
DROP COLUMN rto;

ALTER TABLE business_continuity_plans 
DROP COLUMN mtpd;

ALTER TABLE business_continuity_plans 
DROP COLUMN rpo;
-- +goose StatementEnd
