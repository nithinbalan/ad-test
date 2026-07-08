-- +goose Up
-- +goose StatementBegin
ALTER TABLE vulnerabilities
    ADD qid INT NULL,
        cve NVARCHAR(100) NULL,
        impact NVARCHAR(MAX) NULL,
        solution NVARCHAR(MAX) NULL,
        qualys_severity_id INT  NULL,
        criticality_score INT NULL,
        is_manual BIT NULL DEFAULT 0,
        CONSTRAINT FK_vulnerabilities_qualys_severity_levels  FOREIGN KEY (qualys_severity_id) REFERENCES severity_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @constraintName NVARCHAR(200);

SELECT @constraintName = dc.name
FROM sys.default_constraints dc
JOIN sys.columns c ON dc.parent_object_id = c.object_id AND dc.parent_column_id = c.column_id
WHERE dc.parent_object_id = OBJECT_ID('vulnerabilities') AND c.name = 'is_manual';

IF @constraintName IS NOT NULL
    EXEC('ALTER TABLE vulnerabilities DROP CONSTRAINT ' + @constraintName);

ALTER TABLE vulnerabilities
    DROP CONSTRAINT FK_vulnerabilities_qualys_severity_levels;

ALTER TABLE vulnerabilities
    DROP COLUMN qid, cve, impact, solution, qualys_severity_id, criticality_score, is_manual;
-- +goose StatementEnd