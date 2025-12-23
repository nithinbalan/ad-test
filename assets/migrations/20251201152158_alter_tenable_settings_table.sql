-- +goose Up
-- +goose StatementBegin

ALTER TABLE tenable_settings
ADD status VARCHAR(10) NULL DEFAULT 'active';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

BEGIN
    DECLARE @constraintName NVARCHAR(200)

    SELECT @constraintName = dc.name
    FROM sys.default_constraints dc
    JOIN sys.columns c ON c.default_object_id = dc.object_id
    WHERE c.object_id = OBJECT_ID('dbo.tenable_settings')
      AND c.name = 'status'

    IF @constraintName IS NOT NULL
        EXEC('ALTER TABLE dbo.tenable_settings DROP CONSTRAINT ' + @constraintName)

    ALTER TABLE dbo.tenable_settings DROP COLUMN status
END;

-- +goose StatementEnd
