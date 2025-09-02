-- +goose Up
-- +goose StatementBegin
ALTER TABLE knowledge_center_resources
ALTER COLUMN knowledge_center_resource_status_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE knowledge_center_resources
ALTER COLUMN knowledge_center_resource_status_id INT NOT NULL;
-- +goose StatementEnd
