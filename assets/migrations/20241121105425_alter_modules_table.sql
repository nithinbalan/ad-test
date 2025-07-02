-- +goose Up
-- +goose StatementBegin
ALTER TABLE modules 
ADD is_enable TINYINT NOT NULL DEFAULT 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE modules 
    DROP COLUMN IF EXISTS is_enable;
-- +goose StatementEnd
