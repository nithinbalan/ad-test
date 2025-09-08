-- +goose Up
-- +goose StatementBegin

DECLARE @itemCount INT;

SELECT @itemCount = COUNT(*) FROM strategic_objective_library;

UPDATE knowledge_center_resources
SET item_count = @itemCount
WHERE title_en = 'Strategic Objectives Library';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

UPDATE knowledge_center_resources
SET item_count = 0
WHERE title_en = 'Strategic Objectives Library';

-- +goose StatementEnd
