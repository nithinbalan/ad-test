-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_register 
ADD is_latest TINYINT NOT NULL DEFAULT 1;

ALTER TABLE compliance_register 
ADD compliance_register_id INT;

ALTER TABLE compliance_register
ADD CONSTRAINT FK_compliance_register_parent_compliance_register
FOREIGN KEY (compliance_register_id) REFERENCES compliance_register(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DECLARE @ConstraintName NVARCHAR(200);
SELECT @ConstraintName = dc.name
FROM sys.default_constraints dc
JOIN sys.columns c ON dc.parent_object_id = c.object_id AND dc.parent_column_id = c.column_id
WHERE OBJECT_NAME(dc.parent_object_id) = 'compliance_register' AND c.name = 'is_latest';
IF @ConstraintName IS NOT NULL
BEGIN
    EXEC('ALTER TABLE compliance_register DROP CONSTRAINT ' + @ConstraintName);
END
ALTER TABLE compliance_register DROP COLUMN is_latest;

ALTER TABLE compliance_register
DROP CONSTRAINT FK_compliance_register_parent_compliance_register;

ALTER TABLE compliance_register
DROP COLUMN compliance_register_id;
-- +goose StatementEnd
