-- +goose Up
-- +goose StatementBegin
UPDATE knowledge_center_resources
SET knowledge_center_resource_status_id = NULL
WHERE title_en IN ('KPI Library', 'Strategic Objectives Library');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Get the status ID for 'Coming Soon'
DECLARE @coming_soon_status_id INT;
SELECT @coming_soon_status_id = id FROM knowledge_center_resource_statuses WHERE title_en = 'Coming Soon';

-- Revert statuses for KPI Library and Strategic Objectives Library
UPDATE knowledge_center_resources
SET knowledge_center_resource_status_id = @coming_soon_status_id
WHERE title_en IN ('KPI Library', 'Strategic Objectives Library');

-- +goose StatementEnd
