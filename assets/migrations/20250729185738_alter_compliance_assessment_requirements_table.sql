-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_assessment_requirements
ADD activity NVARCHAR(MAX) NULL;

-- Drop the default constraint on deadline column first
DECLARE @constraint_name NVARCHAR(200)
SELECT @constraint_name = dc.name 
FROM sys.default_constraints dc
INNER JOIN sys.columns c ON dc.parent_object_id = c.object_id AND dc.parent_column_id = c.column_id
WHERE c.name = 'is_completed' AND OBJECT_NAME(dc.parent_object_id) = 'compliance_assessment_requirements'

IF @constraint_name IS NOT NULL
BEGIN
    EXEC('ALTER TABLE compliance_assessment_requirements DROP CONSTRAINT ' + @constraint_name)
END

ALTER TABLE compliance_assessment_requirements
DROP COLUMN is_completed;

ALTER TABLE compliance_assessment_requirements
DROP COLUMN score;

ALTER TABLE compliance_assessment_requirements
DROP COLUMN comments;

ALTER TABLE compliance_assessment_requirements
DROP CONSTRAINT FK_compliance_assessment_requirements_status;

ALTER TABLE compliance_assessment_requirements
DROP COLUMN compliance_assessment_requirement_status_id;

ALTER TABLE compliance_assessment_requirements
DROP CONSTRAINT FK_compliance_assessment_requirements_compliance_statuses

ALTER TABLE compliance_assessment_requirements
DROP COLUMN compliance_status_id;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_assessment_requirements
DROP COLUMN activity;

ALTER TABLE compliance_assessment_requirements
ADD is_completed TINYINT NOT NULL DEFAULT 0;

ALTER TABLE compliance_assessment_requirements
ADD score INT NULL;

ALTER TABLE compliance_assessment_requirements
ADD comments NVARCHAR(MAX) NULL;

ALTER TABLE compliance_assessment_requirements
ADD compliance_status_id INT NULL;

ALTER TABLE compliance_assessment_requirements
ADD compliance_assessment_requirement_status_id INT NULL;

ALTER TABLE compliance_assessment_requirements
ADD CONSTRAINT FK_compliance_assessment_requirements_status FOREIGN KEY (compliance_assessment_requirement_status_id) REFERENCES compliance_assessment_requirement_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE compliance_assessment_requirements
ADD CONSTRAINT FK_compliance_assessment_requirements_compliance_statuses FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd
