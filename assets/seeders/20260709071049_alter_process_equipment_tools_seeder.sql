-- +goose Up
-- +goose StatementBegin
UPDATE process_equipment_tools SET type = 'portable_laptop' WHERE title_en = 'Portable Laptop';
UPDATE process_equipment_tools SET type = 'other' WHERE title_en != 'Portable Laptop';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE process_equipment_tools SET type = 'other';
-- +goose StatementEnd
