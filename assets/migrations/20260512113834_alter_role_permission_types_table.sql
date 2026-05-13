-- +goose Up
-- +goose StatementBegin

DECLARE @constraint_name NVARCHAR(200);

-- Find Existing Constraint
SELECT @constraint_name = cc.name
FROM sys.check_constraints cc
WHERE cc.parent_object_id = OBJECT_ID('role_permission_types')
AND COL_NAME(cc.parent_object_id, cc.parent_column_id) = 'type';

-- Drop Constraint
IF @constraint_name IS NOT NULL
BEGIN
    EXEC('ALTER TABLE role_permission_types DROP CONSTRAINT ' + @constraint_name);
END 

-- Add New Constraint
ALTER TABLE role_permission_types
ADD CONSTRAINT CK_role_permission_types_type
CHECK (type IN ('manage', 'modify', 'view', 'passive_enabled'));

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop New Constraint
ALTER TABLE role_permission_types DROP CONSTRAINT CK_role_permission_types_type;

-- Restore Old Constraint
ALTER TABLE role_permission_types
ADD CONSTRAINT CK_role_permission_types_type
CHECK (type IN ('manage', 'modify', 'view'));

-- +goose StatementEnd
