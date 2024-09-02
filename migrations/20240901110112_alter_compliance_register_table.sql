-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_register ADD total_score INT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_register DROP COLUMN total_score;
-- +goose StatementEnd
