-- +goose Up
-- +goose StatementBegin
-- Insert new types into the document_change_request_types table
INSERT INTO document_change_request_types (title_en, title_ar, type, status)
VALUES
    ('Restore', 'استعادة', 'restore', 'active'),
    ('Restore with Minor Update', 'استعادة مع تحديث طفيف', 'restore-with-minor-update', 'active'),
    ('Restore with Major Update', 'استعادة مع تحديث كبير', 'restore-with-major-update', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the new types from the document_change_request_types table
DELETE FROM document_change_request_types
WHERE type IN ('restore', 'restore-with-minor-update', 'restore-with-major-update');
-- +goose StatementEnd
