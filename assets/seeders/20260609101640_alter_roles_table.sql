-- +goose Up
-- +goose StatementBegin

UPDATE roles
SET
    title_en = 'Committee Admin',
    title_ar = 'مسؤول اللجنة',
    type = 'committee-admin'
WHERE type = 'governance-head';

IF NOT EXISTS (
    SELECT 1
    FROM roles
    WHERE type = 'head-of-hr'
)
BEGIN
    INSERT INTO roles (
        title_en,
        title_ar,
        status,
        is_super_admin,
        license_type_id,
        type
    )
    VALUES (
        'Head of HR',
        'رئيس الموارد البشرية',
        'active',
        0,
        1,
        'head-of-hr'
    );
END

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DELETE FROM roles
WHERE type = 'head-of-hr';

UPDATE roles
SET
    title_en = 'Governance Head',
    title_ar = 'رئيس الحوكمة',
    type = 'governance-head'
WHERE type = 'committee-admin';

-- +goose StatementEnd