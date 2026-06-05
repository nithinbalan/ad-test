
-- +goose Up
-- +goose StatementBegin
ALTER TABLE workflow_level_types
ADD module_identifier NVARCHAR(255) NOT NULL DEFAULT 'common';
-- +goose StatementEnd




-- +goose Down
-- +goose StatementBegin
DECLARE @ConstraintName NVARCHAR(128);

SELECT @ConstraintName = name
FROM sys.default_constraints
WHERE parent_object_id = OBJECT_ID('workflow_level_types')
  AND COL_NAME(parent_object_id, parent_column_id) = 'module_identifier';

IF @ConstraintName IS NOT NULL
    EXEC('ALTER TABLE workflow_level_types DROP CONSTRAINT ' + @ConstraintName);

ALTER TABLE workflow_level_types
DROP COLUMN module_identifier;
-- +goose StatementEnd