-- +goose Up
-- +goose StatementBegin
DECLARE @UserID INT;
DECLARE @AdminDesignationID INT;

SELECT @AdminDesignationID = id FROM designations WHERE title_en = 'Admin';

-- Insert a user
INSERT INTO users (
    organization_unit_1_id,
    organization_unit_2_id,
    first_name_en,
    last_name_en,
    first_name_ar,
    last_name_ar,
    email,
    user_status_id,
    designation_id,
    profile_image_path,
    initial,
    color,
    created_by_id,
    created_at,
    updated_at
) VALUES (
    (SELECT id FROM organization_units WHERE title_en = 'Information Technology'),
    (SELECT id FROM organization_units WHERE title_en = 'IT Support'),
    'Super',
    'Admin',
    'سوبر',
    'أدمن',
    'admin@jethur.com',
    1,
    @AdminDesignationID,
    NULL,
    'SA',
    '#FFD700',
    1,
    GETDATE(),
    GETDATE()
);

-- Get the ID of the inserted user
SET @UserID = SCOPE_IDENTITY();

-- Insert a corresponding user_security record with the bcrypt hash
INSERT INTO user_securities (
    user_id,
    password_hash,
    is_mfa,
    otp,
    otp_expiry_at,
    updated_at
) VALUES (
    @UserID,
    '$2a$10$lD/krN00lAIwGgy4gBVcPuG2lnY6xIdM8Xq74.eDgPDQ8g/rj8scq', -- Replace with your actual bcrypt hash
    0,
    NULL,
    NULL,
    GETDATE()
);

-- Insert a corresponding user_roles record with the bcrypt hash
INSERT INTO user_roles (
    user_id,
    role_id
) VALUES (
    @UserID,
    1
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Here, add your logic to clean up the data for the Down migration
DELETE FROM user_roles WHERE user_id = @UserID;
DELETE FROM user_securities WHERE user_id = @UserID;
DELETE FROM users WHERE id = @UserID;
-- +goose StatementEnd
