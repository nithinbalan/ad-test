-- +goose Up
-- +goose StatementBegin
DECLARE @UserID INT;
DECLARE @AdminDesignationID INT;
DECLARE @SuperAdminRoleID INT;

-- Get Admin designation and Super Admin role
SELECT @AdminDesignationID = id FROM designations WHERE title_en = 'Admin';
SELECT @SuperAdminRoleID = id FROM roles WHERE is_super_admin = 1;

-- Insert Admin user
INSERT INTO users (first_name_en, last_name_en, first_name_ar, last_name_ar, email, user_status_id, designation_id, profile_image_path, initial, color, created_by_id, created_at, updated_at) 
VALUES ( 'System', '', 'النظام', '', 'no-reply@jethur.com', 1, @AdminDesignationID, NULL, 'S', '#4A90E2', 1, GETDATE(), GETDATE());

-- Get the inserted user ID
SET @UserID = SCOPE_IDENTITY();

-- Insert user security
INSERT INTO user_securities (user_id, password_hash, is_mfa, otp, otp_expiry_at, updated_at) 
VALUES (@UserID, '$2a$12$1ei2lpQgem/SH3jUUYd2h.1uJToEcQcTW3diTuBsxuudR4.lV.9Jm', 0, NULL, NULL, GETDATE());

-- Assign Super Admin role
INSERT INTO user_roles (user_id, role_id)
VALUES (@UserID, @SuperAdminRoleID);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @UserID INT;

SELECT @UserID = id FROM users WHERE email = 'no-reply@jethur.com';

DELETE FROM user_roles WHERE user_id = @UserID;
DELETE FROM user_securities WHERE user_id = @UserID;
DELETE FROM users WHERE id = @UserID;
-- +goose StatementEnd
