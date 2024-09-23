-- +goose Up
-- +goose StatementBegin
ALTER TABLE audit_finding_categories
ADD color_code VARCHAR(10) NULL; 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE audit_finding_categories
DROP COLUMN color_code;
-- +goose StatementEnd
