-- +goose Up
-- +goose StatementBegin
DECLARE @PasswordHash VARBINARY(MAX);
SET @PasswordHash = HASHBYTES('SHA2_256', '12345678');

INSERT INTO users (
    organization_unit_1_id,
    organization_unit_2_id,
    first_name_en,
    last_name_en,
    first_name_ar,
    last_name_ar,
    email,
    password_hash,
    profile_image_path,
    created_at,
    updated_at
) VALUES (
    (SELECT id FROM organization_units WHERE title_en = 'Information Technology Division'),
    (SELECT id FROM organization_units WHERE title_en = 'IT Support Department'),
    'Super',
    'Admin',
    'سوبر',
    'أدمن',
    'admin@jethur.com',
    CAST('' AS XML).value('xs:base64Binary(sql:variable("@PasswordHash"))', 'VARCHAR(MAX)'),
    NULL,
    GETDATE(),
    GETDATE()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE users;
-- +goose StatementEnd
