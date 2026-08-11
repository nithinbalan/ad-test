-- +goose Up
-- +goose StatementBegin
INSERT INTO incident_impacts (title_en, title_ar, description, type, status, created_by_id, created_at)
VALUES
    ('Confidentiality', N'السرية', NULL, 'confidentiality', 'active', 1, GETDATE()),
    ('Integrity', N'النزاهة', NULL, 'integrity', 'active', 1, GETDATE()),
    ('Availability', N'التوافر', NULL, 'availability', 'active', 1, GETDATE());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM incident_impacts
WHERE title_en IN (
    'Confidentiality',
    'Integrity',
    'Availability'
);
-- +goose StatementEnd
