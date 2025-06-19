-- +goose Up
-- +goose StatementBegin

-- Drop existing foreign key constraints
ALTER TABLE processes DROP CONSTRAINT FK_processes_users_1;
ALTER TABLE processes DROP CONSTRAINT FK_processes_organization_units;

-- Rename the columns
EXEC sp_rename 'processes.owner_id', 'responsible_user_id', 'COLUMN';
EXEC sp_rename 'processes.owner_unit_id', 'responsible_unit_id', 'COLUMN';

-- Re-add foreign keys with updated column names
ALTER TABLE processes
ADD CONSTRAINT FK_processes_responsible_user FOREIGN KEY (responsible_user_id)
REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE processes
ADD CONSTRAINT FK_processes_responsible_unit FOREIGN KEY (responsible_unit_id)
REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop new foreign key constraints
ALTER TABLE processes DROP CONSTRAINT FK_processes_responsible_user;
ALTER TABLE processes DROP CONSTRAINT FK_processes_responsible_unit;

-- Rename columns back to original
EXEC sp_rename 'processes.responsible_user_id', 'owner_id', 'COLUMN';
EXEC sp_rename 'processes.responsible_unit_id', 'owner_unit_id', 'COLUMN';

-- Re-add original foreign keys
ALTER TABLE processes
ADD CONSTRAINT FK_processes_users_1 FOREIGN KEY (owner_id)
REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE processes
ADD CONSTRAINT FK_processes_organization_units FOREIGN KEY (owner_unit_id)
REFERENCES organization_units(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd
