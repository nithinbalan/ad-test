-- +goose Up
-- +goose StatementBegin

UPDATE knowledge_center_resource_statuses
SET type = 'new'
WHERE title_en = 'New';

UPDATE knowledge_center_resource_statuses
SET type = 'coming-soon'
WHERE title_en = 'Coming Soon';

UPDATE knowledge_center_resource_statuses
SET type = 'production-only'
WHERE title_en = 'Production';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

UPDATE knowledge_center_resource_statuses
SET type = NULL;

-- +goose StatementEnd
