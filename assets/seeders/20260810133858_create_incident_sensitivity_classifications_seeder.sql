-- +goose Up
-- +goose StatementBegin
INSERT INTO incident_sensitivity_classifications (title_en, title_ar, description, type, status, created_by_id, created_at)
VALUES
    ('Public', N'عام', NULL, 'public', 'active', 1, GETDATE()),
    ('Restricted', N'مقيد', NULL, 'restricted', 'active', 1, GETDATE()),
    ('Confidential', N'سري', NULL, 'confidential', 'active', 1, GETDATE());
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM incident_sensitivity_classifications
WHERE type IN (
    'public',
    'restricted',
    'confidential'
);
-- +goose StatementEnd
