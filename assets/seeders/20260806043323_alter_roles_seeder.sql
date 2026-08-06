-- +goose Up
-- +goose StatementBegin

INSERT INTO roles (
    title_en,
    title_ar,
    status,
    is_super_admin,
    license_type_id,
    type
)
VALUES
(
    'CEO',
    'الرئيس التنفيذي',
    'active',
    0,
    1,
    'ceo'
),
(
    'Shared Services',
    'الخدمات المشتركة',
    'active',
    0,
    1,
    'shared-services'
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DELETE FROM roles
WHERE type IN ('ceo', 'shared-services');

-- +goose StatementEnd