-- +goose Up
-- +goose StatementBegin

UPDATE audit_statuses
SET color_code = '#007BFF'
WHERE type = 'planning';

UPDATE audit_statuses
SET color_code = '#FF9800'
WHERE type = 'fieldwork';

UPDATE audit_statuses
SET color_code = '#D6A2E8'
WHERE type = 'reporting';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

UPDATE audit_statuses
SET color_code = '#3B82F6'
WHERE type = 'planning';

UPDATE audit_statuses
SET color_code = '#F59E0B'
WHERE type = 'fieldwork';

UPDATE audit_statuses
SET color_code = '#10B981'
WHERE type = 'reporting';

-- +goose StatementEnd