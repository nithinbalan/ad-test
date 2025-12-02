-- +goose Up
-- +goose StatementBegin

UPDATE integrations
SET 
    title_en = 'Microsoft Entra ID',
    title_ar = 'مايكروسوفت إنتراإي دي',
    description_en = 'Microsoft Entra ID integration for authentication and directory services.',
    description_ar = 'تكامل Microsoft Entra ID للمصادقة وخدمات الدليل.'
WHERE title_en = 'Azure AD';

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

UPDATE integrations
SET 
    title_en = 'Azure AD',
    title_ar = 'أزور آي دي',
    description_en = 'Microsoft Azure Active Directory integration for authentication and directory services.',
    description_ar = 'تكامل Microsoft Azure Active Directory للمصادقة وخدمات الدليل.'
WHERE title_en = 'Microsoft Entra ID';

-- +goose StatementEnd
