-- +goose Up
-- +goose StatementBegin

ALTER TABLE knowledge_center_resources
ADD type VARCHAR(50) NULL;

ALTER TABLE knowledge_center_resources
ADD item_count INT NOT NULL DEFAULT 0;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE knowledge_center_resources
DROP COLUMN type;

ALTER TABLE knowledge_center_resources
DROP COLUMN item_count;

-- +goose StatementEnd
