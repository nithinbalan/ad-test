-- +goose Up
-- +goose StatementBegin

ALTER TABLE third_party_assessments
ADD is_send TINYINT  NOT NULL DEFAULT 0;

ALTER TABLE third_party_assessments
ADD imported_by_id INT NULL;

ALTER TABLE third_party_assessments
ADD imported_at DATETIME NULL;

ALTER TABLE third_party_assessments
ADD CONSTRAINT FK_third_party_assessments_import_user FOREIGN KEY (imported_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE third_party_assessments
ADD score FLOAT NULL;

ALTER TABLE third_party_assessments
ADD risk_level_id INT NULL;

ALTER TABLE third_party_assessments
ADD CONSTRAINT FK_third_party_assessments_risk_levels
FOREIGN KEY (risk_level_id) REFERENCES risk_levels(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE third_party_assessments 
ADD sent_by_id INT;

ALTER TABLE third_party_assessments
ADD CONSTRAINT FK_third_party_assessments_sent_user
FOREIGN KEY (sent_by_id) REFERENCES users(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE third_party_assessments 
ADD sent_at DATETIME NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DECLARE @ConstraintName NVARCHAR(200);
SELECT @ConstraintName = dc.name
FROM sys.default_constraints dc
JOIN sys.columns c ON dc.parent_object_id = c.object_id AND dc.parent_column_id = c.column_id
WHERE OBJECT_NAME(dc.parent_object_id) = 'third_party_assessments'
  AND c.name = 'is_send';

IF @ConstraintName IS NOT NULL
BEGIN
    EXEC('ALTER TABLE third_party_assessments DROP CONSTRAINT ' + @ConstraintName);
END

-- Now drop the column
ALTER TABLE third_party_assessments DROP COLUMN is_send;

ALTER TABLE third_party_assessments
DROP CONSTRAINT FK_third_party_assessments_import_user;

ALTER TABLE third_party_assessments
DROP COLUMN imported_by_id;

ALTER TABLE third_party_assessments
DROP COLUMN imported_at;

ALTER TABLE third_party_assessments
DROP CONSTRAINT FK_third_party_assessments_sent_user;

ALTER TABLE third_party_assessments
DROP COLUMN sent_by_id;

ALTER TABLE third_party_assessments
DROP COLUMN sent_at;

ALTER TABLE third_party_assessments
DROP CONSTRAINT FK_third_party_assessments_risk_levels;


ALTER TABLE third_party_assessments
DROP COLUMN risk_level_id;

ALTER TABLE third_party_assessments
DROP COLUMN score;
-- +goose StatementEnd
