-- +goose Up
-- +goose StatementBegin
-- Add types to existing seeded roles
UPDATE roles SET type = 'super-admin' WHERE title_en = 'Super Admin';
UPDATE roles SET type = 'auditor' WHERE title_en = 'Auditor';
UPDATE roles SET type = 'risk-manager' WHERE title_en = 'Risk Manager';
UPDATE roles SET type = 'risk-executive' WHERE title_en = 'Risk Executive';

UPDATE roles SET type = 'other' where title_en NOT IN('Super Admin','Auditor','Risk Manager','Risk Executive')

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE roles SET type = NULL;
-- +goose StatementEnd
