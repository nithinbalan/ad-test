-- +goose Up
-- +goose StatementBegin

DECLARE @superAdminRoleID INT;
SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;

INSERT INTO risk_settings (review_role_id, created_by_id, created_at, updated_by_id, updated_at)
VALUES (@superAdminRoleID, 1, GETDATE(), NULL, NULL);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE risk_settings;
-- +goose StatementEnd
