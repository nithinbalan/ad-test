-- +goose Up
-- +goose StatementBegin
INSERT INTO roles (
    title_en,
    title_ar,
    status,
    is_super_admin,
    license_type_id,
    type
) VALUES (
    'Governance Head',
    'رئيس الحوكمة',
    'active',
    0,
    1,
    'governance-head'
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM roles
WHERE type = 'governance-head';
-- +goose StatementEnd