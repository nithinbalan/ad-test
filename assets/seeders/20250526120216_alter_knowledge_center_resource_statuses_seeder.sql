-- +goose Up
-- +goose StatementBegin

UPDATE knowledge_center_resource_statuses
SET color_code = '#28A745'
WHERE title_en = 'New';

UPDATE knowledge_center_resource_statuses
SET color_code = '#FFC107'
WHERE title_en = 'Coming Soon';

UPDATE knowledge_center_resource_statuses
SET color_code = '#007BFF'
WHERE title_en = 'Production';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

UPDATE knowledge_center_resource_statuses
SET color_code = NULL;

-- +goose StatementEnd
