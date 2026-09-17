-- +goose Up
-- +goose StatementBegin
UPDATE user_securities 
SET password_hash = '$2a$10$c.b1m.D.d/lzSUwbqykyo.09ztmzrHuZESoTOJ5itnCEjEJ6mu7Na',
    updated_at = GETDATE()
WHERE user_id = (SELECT id FROM users WHERE email = 'globaladmin@jethur.com');

UPDATE user_securities
SET password_hash = '$2a$10$l3p2ldDXKDCWkRPD4vzvuu3bneGgNzDn4jG6L8aYxEimrUxIw6Y8u',
    updated_at = GETDATE()
WHERE user_id = (SELECT id FROM users WHERE email = 'no-reply@jethur.com');

-- Lock the updated users by resolving the 'locked' user_status type
UPDATE users
SET user_status_id = (SELECT id FROM user_statuses WHERE type = 'locked')
WHERE email IN ('globaladmin@jethur.com', 'no-reply@jethur.com');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE user_securities
SET password_hash = '$2a$10$y8CwihVOmMS6pJNnKMm4quofaa2eo1RvaVIOpBpYNE4fy0MkWnVL2',
    updated_at = GETDATE()
WHERE user_id = (SELECT id FROM users WHERE email = 'globaladmin@jethur.com');

UPDATE user_securities
SET password_hash = '$2a$12$1ei2lpQgem/SH3jUUYd2h.1uJToEcQcTW3diTuBsxuudR4.lV.9Jm',
    updated_at = GETDATE()
WHERE user_id = (SELECT id FROM users WHERE email = 'no-reply@jethur.com');

-- Restore the users to the 'active' status
UPDATE users
SET user_status_id = (SELECT id FROM user_statuses WHERE type = 'active')
WHERE email IN ('globaladmin@jethur.com', 'no-reply@jethur.com');
-- +goose StatementEnd
