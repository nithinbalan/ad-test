-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_statuses
DROP CONSTRAINT CHK_document_statuses_type;

ALTER TABLE document_statuses
ADD CONSTRAINT CHK_document_statuses_type CHECK (type IN ('draft', 'in-review', 'rejected', 'approved', 'published', 'archived', 'reverted'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE document_statuses
DROP CONSTRAINT CHK_document_statuses_type;

ALTER TABLE document_statuses
ADD CONSTRAINT CHK_document_statuses_type CHECK (type IN ('draft', 'in-review', 'rejected', 'approved', 'published', 'archived'));
-- +goose StatementEnd
