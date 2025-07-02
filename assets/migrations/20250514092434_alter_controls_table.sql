-- +goose Up
-- +goose StatementBegin

-- Drop existing foreign key constraints
ALTER TABLE controls DROP CONSTRAINT FK_controls_user_1;
ALTER TABLE controls DROP CONSTRAINT FK_controls_organization_units;

-- Rename the columns
EXEC sp_rename 'controls.owner_id', 'responsible_user_id', 'COLUMN';
EXEC sp_rename 'controls.owner_unit_id', 'responsible_unit_id', 'COLUMN';

-- Re-add foreign key constraints with new column names
ALTER TABLE controls
ADD CONSTRAINT FK_controls_responsible_user FOREIGN KEY (responsible_user_id)
REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE controls
ADD CONSTRAINT FK_controls_responsible_unit FOREIGN KEY (responsible_unit_id)
REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop newly added foreign keys
ALTER TABLE controls DROP CONSTRAINT FK_controls_responsible_user;
ALTER TABLE controls DROP CONSTRAINT FK_controls_responsible_unit;

-- Rename columns back to original
EXEC sp_rename 'controls.responsible_user_id', 'owner_id', 'COLUMN';
EXEC sp_rename 'controls.responsible_unit_id', 'owner_unit_id', 'COLUMN';

-- Re-add original foreign keys
ALTER TABLE controls
ADD CONSTRAINT FK_controls_user_1 FOREIGN KEY (owner_id)
REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE controls
ADD CONSTRAINT FK_controls_organization_units FOREIGN KEY (owner_unit_id)
REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd
