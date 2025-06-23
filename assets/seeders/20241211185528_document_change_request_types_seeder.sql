-- +goose Up
-- +goose StatementBegin
INSERT INTO document_change_request_types (title_en, title_ar, type, status)
VALUES
    ('Minor Update', 'تحديث طفيف', 'minor', 'active'),
    ('Major Update', 'تحديث كبير', 'major', 'active'),
    ('Archive', 'أرشفة', 'archive', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM document_change_request_types
WHERE type IN ('minor', 'major', 'archive');
-- +goose StatementEnd
