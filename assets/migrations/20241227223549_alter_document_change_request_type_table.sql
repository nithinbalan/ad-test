-- +goose Up
-- +goose StatementBegin
-- Drop the existing CHECK constraint on the 'type' column
ALTER TABLE document_change_request_types DROP CONSTRAINT CHK_document_change_request_types;

-- Re-add the CHECK constraint with the new 'restore' type included
ALTER TABLE document_change_request_types
ADD CONSTRAINT CHK_document_change_request_types CHECK (type IN ('minor', 'major', 'archive', 'restore', 'restore-with-minor-update', 'restore-with-major-update'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop the updated CHECK constraint
ALTER TABLE document_change_request_types DROP CONSTRAINT CHK_document_change_request_types;

-- Re-add the original CHECK constraint without 'restore'
ALTER TABLE document_change_request_types
ADD CONSTRAINT CHK_document_change_request_types CHECK (type IN ('minor', 'major', 'archive'));
-- +goose StatementEnd
