-- +goose Up
-- +goose StatementBegin
ALTER TABLE escalations
ADD [type] NVARCHAR(MAX) NOT NULL
    CONSTRAINT [DF_escalations_type] DEFAULT ('type') WITH VALUES;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @constraint_name NVARCHAR(128);

SELECT @constraint_name = dc.name
FROM sys.default_constraints dc
JOIN sys.columns c
    ON c.default_object_id = dc.object_id
JOIN sys.tables t
    ON t.object_id = c.object_id
WHERE t.name = 'escalations'
  AND c.name = 'type';

IF @constraint_name IS NOT NULL
    EXEC ('ALTER TABLE escalations DROP CONSTRAINT [' + @constraint_name + ']');

ALTER TABLE escalations
DROP COLUMN [type];
-- +goose StatementEnd
