-- +goose Up
-- +goose StatementBegin
INSERT INTO document_change_request_types (title_en, title_ar, type, status)
VALUES
    ('Modify Document', 'تعديل الوثيقة', 'modify', 'active'),
    ('Archive Document', 'أرشفة الوثيقة', 'archive', 'active'),
    ('Restore Document', 'استعادة الوثيقة', 'restore', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM document_change_request_types 
WHERE type IN ('modify', 'archive', 'restore');
-- +goose StatementEnd