-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_versions
ADD acknowledge_due_date DATE NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE document_versions
DROP COLUMN acknowledge_due_date;
-- +goose StatementEnd
