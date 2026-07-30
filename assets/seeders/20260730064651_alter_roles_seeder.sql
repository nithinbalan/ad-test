-- +goose Up
-- +goose StatementBegin

DECLARE @type_id INT;
SELECT @type_id = id
FROM license_types
WHERE type = 'active';

INSERT INTO roles (
    title_en,
    title_ar,
    status,
    is_super_admin,
    license_type_id,
    type
)
VALUES (
    'Cybersecurity Supervisory Committee',
    N'لجنة الإشراف على الأمن السيبراني',
    'active',
    0,
    @type_id,
    'cybersecurity-supervisory-committee'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM roles
WHERE type = 'cybersecurity-supervisory-committee';

-- +goose StatementEnd