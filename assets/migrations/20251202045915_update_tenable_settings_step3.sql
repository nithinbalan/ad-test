-- +goose Up
-- +goose NO TRANSACTION

BEGIN
    DECLARE @constraintName NVARCHAR(200)

    SELECT @constraintName = dc.name
    FROM sys.default_constraints dc
    JOIN sys.columns c ON c.default_object_id = dc.object_id
    WHERE c.object_id = OBJECT_ID('tenable_settings')
      AND c.name = 'is_active'

    IF @constraintName IS NOT NULL
        EXEC('ALTER TABLE tenable_settings DROP CONSTRAINT ' + @constraintName)

    ALTER TABLE tenable_settings DROP COLUMN is_active
END;
