-- +goose Up
-- +goose StatementBegin
-- Add the recovery objective columns to the asset_register table
ALTER TABLE asset_register 
ADD rto INT NULL;

ALTER TABLE asset_register 
ADD rto_unit_id INT NULL;

ALTER TABLE asset_register 
ADD mtpd INT NULL;

ALTER TABLE asset_register 
ADD mtpd_unit_id INT NULL;

ALTER TABLE asset_register 
ADD rpo INT NULL;

ALTER TABLE asset_register 
ADD rpo_unit_id INT NULL;

-- Add the foreign key constraints for rto_unit_id, mtpd_unit_id, and rpo_unit_id
ALTER TABLE asset_register
ADD CONSTRAINT FK_asset_register_rto_time_units FOREIGN KEY (rto_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE asset_register
ADD CONSTRAINT FK_asset_register_mtpd_time_units FOREIGN KEY (mtpd_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE asset_register
ADD CONSTRAINT FK_asset_register_rpo_time_units FOREIGN KEY (rpo_unit_id) REFERENCES time_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove foreign key constraints before dropping columns
ALTER TABLE asset_register
DROP CONSTRAINT FK_asset_register_rto_time_units;

ALTER TABLE asset_register
DROP COLUMN rto_unit_id;

ALTER TABLE asset_register
DROP CONSTRAINT FK_asset_register_mtpd_time_units;

ALTER TABLE asset_register
DROP COLUMN mtpd_unit_id;

ALTER TABLE asset_register
DROP CONSTRAINT FK_asset_register_rpo_time_units;

ALTER TABLE asset_register
DROP COLUMN rpo_unit_id;

-- Remove the additional columns
ALTER TABLE asset_register 
DROP COLUMN rto;

ALTER TABLE asset_register 
DROP COLUMN mtpd;

ALTER TABLE asset_register 
DROP COLUMN rpo;
-- +goose StatementEnd
