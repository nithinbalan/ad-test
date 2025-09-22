-- +goose Up
-- +goose StatementBegin
DECLARE @dfname SYSNAME, @schema SYSNAME, @sql NVARCHAR(MAX);

SELECT 
  @dfname = dc.name,
  @schema = s.name
FROM sys.default_constraints AS dc
JOIN sys.columns AS c
  ON dc.parent_object_id = c.object_id
 AND dc.parent_column_id = c.column_id
JOIN sys.tables AS t
  ON t.object_id = dc.parent_object_id
JOIN sys.schemas AS s
  ON s.schema_id = t.schema_id
WHERE t.name = 'meeting_action_plans'
  AND c.name = 'review_user_id';

IF @dfname IS NOT NULL
BEGIN
    SET @sql = N'ALTER TABLE ' 
             + QUOTENAME(@schema) + N'.' + QUOTENAME(N'meeting_action_plans')
             + N' DROP CONSTRAINT ' + QUOTENAME(@dfname) + N';';
    EXEC sp_executesql @sql;
END
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @schema SYSNAME, @sql NVARCHAR(MAX);

SELECT @schema = s.name
FROM sys.tables AS t
JOIN sys.schemas AS s ON s.schema_id = t.schema_id
WHERE t.name = 'meeting_action_plans';

IF NOT EXISTS (
    SELECT 1
    FROM sys.default_constraints AS dc
    JOIN sys.columns AS c
      ON dc.parent_object_id = c.object_id
     AND dc.parent_column_id = c.column_id
    JOIN sys.tables AS t
      ON t.object_id = dc.parent_object_id
    WHERE t.name = 'meeting_action_plans'
      AND c.name = 'review_user_id'
)
BEGIN
    SET @sql = N'ALTER TABLE '
             + QUOTENAME(@schema) + N'.' + QUOTENAME(N'meeting_action_plans')
             + N' ADD CONSTRAINT DF_meeting_action_plans_review_user_id DEFAULT (0) FOR review_user_id;';
    EXEC sp_executesql @sql;
END
-- +goose StatementEnd
