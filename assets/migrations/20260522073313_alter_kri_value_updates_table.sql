-- +goose Up
-- +goose StatementBegin

ALTER TABLE kri_value_updates
ADD is_previous BIT NOT NULL
    CONSTRAINT DF_kri_value_updates_is_previous DEFAULT 0;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE kri_value_updates
DROP CONSTRAINT DF_kri_value_updates_is_previous;

ALTER TABLE kri_value_updates
DROP COLUMN is_previous;

-- +goose StatementEnd