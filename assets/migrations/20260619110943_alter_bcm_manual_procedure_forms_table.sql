-- +goose Up
-- +goose StatementBegin
-- Drop existing foreign key constraints
ALTER TABLE bcm_manual_procedure_forms
DROP CONSTRAINT FK_bcm_manual_procedure_forms_users_1;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE bcm_manual_procedure_forms
DROP CONSTRAINT FK_bcm_manual_procedure_forms_stakeholders_1;
-- +goose StatementEnd

-- +goose StatementBegin
-- Rename owner_id to owner_entity_id for clarity
EXEC sp_rename 'bcm_manual_procedure_forms.owner_id', 'owner_entity_id', 'COLUMN';
-- +goose StatementEnd

-- +goose StatementBegin
-- Add new foreign key to organization_units
ALTER TABLE bcm_manual_procedure_forms
ADD CONSTRAINT FK_bcm_manual_procedure_forms_organization_units_1 
FOREIGN KEY (owner_entity_id) REFERENCES organization_units(id) 
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose StatementBegin
-- Add temporary column for internal stakeholder text
ALTER TABLE bcm_manual_procedure_forms
ADD internal_stakeholder NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose StatementBegin
-- Drop old stakeholder_id column
ALTER TABLE bcm_manual_procedure_forms
DROP COLUMN stakeholder_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Add back stakeholder_id as INT
ALTER TABLE bcm_manual_procedure_forms
ADD stakeholder_id INT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
-- Drop internal_stakeholder column
ALTER TABLE bcm_manual_procedure_forms
DROP COLUMN internal_stakeholder;
-- +goose StatementEnd

-- +goose StatementBegin
-- Drop organization_units foreign key
ALTER TABLE bcm_manual_procedure_forms
DROP CONSTRAINT FK_bcm_manual_procedure_forms_organization_units_1;
-- +goose StatementEnd

-- +goose StatementBegin
-- Rename owner_entity_id back to owner_id
EXEC sp_rename 'bcm_manual_procedure_forms.owner_entity_id', 'owner_id', 'COLUMN';
-- +goose StatementEnd

-- +goose StatementBegin
-- Add back foreign key to users
ALTER TABLE bcm_manual_procedure_forms
ADD CONSTRAINT FK_bcm_manual_procedure_forms_users_1 
FOREIGN KEY (owner_id) REFERENCES users(id) 
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose StatementBegin
-- Add back foreign key to stakeholders
ALTER TABLE bcm_manual_procedure_forms
ADD CONSTRAINT FK_bcm_manual_procedure_forms_stakeholders_1 
FOREIGN KEY (stakeholder_id) REFERENCES stakeholders(id) 
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd
