-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_statuses ADD base_score INT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_statuses DROP COLUMN base_score;
-- +goose StatementEnd
