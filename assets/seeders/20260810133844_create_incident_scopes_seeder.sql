-- +goose Up
-- +goose StatementBegin
INSERT INTO incident_scopes (title_en, title_ar, description, type, status, created_by_id, created_at)
VALUES
    ('Single System', N'نظام واحد', NULL, 'single_system', 'active', 1, GETDATE()),
    ('Multiple Systems', N'أنظمة متعددة', NULL, 'multiple_systems', 'active', 1, GETDATE()),
    ('Whole Organization', N'المنظمة بأكملها', NULL, 'whole_organization', 'active', 1, GETDATE());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM incident_scopes
WHERE type IN (
    'single_system',
    'multiple_systems',
    'whole_organization'
);
-- +goose StatementEnd
