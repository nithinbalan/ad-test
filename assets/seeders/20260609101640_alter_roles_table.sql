-- +goose Up
-- +goose StatementBegin

DELETE FROM roles
WHERE type = 'governance-head';

INSERT INTO roles (
    title_en,
    title_ar,
    status,
    is_super_admin,
    license_type_id,
    type
) VALUES
(
    'Committee Admin',
    'مسؤول اللجنة',
    'active',
    0,
    1,
    'committee-admin'
),
(
    'Head of HR',
    'رئيس الموارد البشرية',
    'active',
    0,
    1,
    'head-of-hr'
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DELETE FROM roles
WHERE type IN ('committee-admin', 'head-of-hr');

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