-- +goose Up
-- +goose StatementBegin

ALTER TABLE compliance_register
ADD parsed_content NVARCHAR(MAX) NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE compliance_register
DROP COLUMN parsed_content;

-- +goose StatementEnd
