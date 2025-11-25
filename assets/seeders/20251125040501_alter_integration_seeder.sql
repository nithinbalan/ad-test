-- +goose Up
-- +goose StatementBegin

-- Insert new Azure AD integration using category lookup
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
    (SELECT id FROM integration_categories WHERE title_en = 'Single Sign-On (SSO)'),
    'Azure AD',
    'أزور آي دي',
    'Microsoft Azure Active Directory integration for authentication and directory services.',
    'تكامل Microsoft Azure Active Directory للمصادقة وخدمات الدليل.',
    'azure-ad',
    NULL,
    0,
    'azure-ad',
    GETDATE(),
    1,
    'active'
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

-- Rollback: delete the inserted Azure AD integration
DELETE FROM integrations
WHERE title_en = 'Azure AD'
  AND integration_category_id = (SELECT id FROM integration_categories WHERE title_en = 'Single Sign-On (SSO)');

-- +goose StatementEnd
