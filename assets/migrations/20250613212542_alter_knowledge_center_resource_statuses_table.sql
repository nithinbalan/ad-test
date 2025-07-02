-- +goose Up
-- +goose StatementBegin
ALTER TABLE knowledge_center_resource_statuses
ADD type VARCHAR(50) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE knowledge_center_resource_statuses
DROP COLUMN type;
-- +goose StatementEnd
