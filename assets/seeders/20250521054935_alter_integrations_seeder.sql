-- +goose Up
-- +goose StatementBegin

-- Insert new SIEM integration using category lookup
INSERT INTO integrations (
    integration_category_id,
    title_en,
    title_ar,
    description_en,
    description_ar,
    logo_url,
    api_base_url,
    is_enable,
    type,
    created_at,
    created_by_id,
    status
) VALUES (
    (SELECT id FROM integration_categories WHERE title_en = 'SIEM'),
    'SIEM',
    'سييم',
    'Security Information and Event Management integration for forwarding logs to external systems.',
    'تكامل إدارة معلومات وأحداث الأمان لإرسال السجلات إلى الأنظمة الخارجية.',
    'siem',
    NULL,
    0,
    'siem',
    GETDATE(),
    1,
    'active'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Rollback: delete the inserted SIEM integration
DELETE FROM integrations
WHERE title_en = 'SIEM'
  AND integration_category_id = (SELECT id FROM integration_categories WHERE title_en = 'SIEM');

-- +goose StatementEnd
