-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_requirements ADD score INT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_requirements DROP COLUMN score;
-- +goose StatementEnd
