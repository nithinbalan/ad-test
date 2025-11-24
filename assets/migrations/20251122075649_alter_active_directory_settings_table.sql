-- +goose Up

-- +goose StatementBegin
ALTER TABLE active_directory_settings 
ADD tenant_id NVARCHAR(100) NULL;

ALTER TABLE active_directory_settings 
ADD client_id NVARCHAR(200) NULL;

ALTER TABLE active_directory_settings 
ADD client_secret NVARCHAR(500) NULL;

ALTER TABLE active_directory_settings 
ADD redirect_uri NVARCHAR(500) NULL;

ALTER TABLE active_directory_settings 
ADD auth_type NVARCHAR(50) NOT NULL DEFAULT 'ldap';

ALTER TABLE active_directory_settings 
ADD is_enabled BIT NOT NULL DEFAULT 1;
-- +goose StatementEnd


-- +goose StatementBegin
-- Run in a separate batch so SQL Server can now "see" the column
ALTER TABLE active_directory_settings
ADD CONSTRAINT CK_active_directory_settings_auth_type
CHECK (auth_type IN ('ldap', 'azure'));
-- +goose StatementEnd


-- +goose Down

-- +goose StatementBegin
DECLARE @ConstraintName NVARCHAR(200);
SELECT @ConstraintName = dc.name
FROM sys.check_constraints dc
WHERE OBJECT_NAME(dc.parent_object_id) = 'active_directory_settings'
  AND dc.name LIKE 'CK_active_directory_settings_auth_type%';
IF @ConstraintName IS NOT NULL
BEGIN
    EXEC('ALTER TABLE active_directory_settings DROP CONSTRAINT ' + @ConstraintName);
END
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE active_directory_settings DROP COLUMN tenant_id;
ALTER TABLE active_directory_settings DROP COLUMN client_id;
ALTER TABLE active_directory_settings DROP COLUMN client_secret;
ALTER TABLE active_directory_settings DROP COLUMN redirect_uri;
ALTER TABLE active_directory_settings DROP COLUMN auth_type;
ALTER TABLE active_directory_settings DROP COLUMN is_enabled;
-- +goose StatementEnd
