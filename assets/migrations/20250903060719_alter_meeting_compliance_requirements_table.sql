-- +goose Up
-- +goose StatementBegin
ALTER TABLE meeting_compliance_requirements
ALTER COLUMN compliance_requirement_id INT NULL;

-- Add nullable compliance_clause_id column
ALTER TABLE meeting_compliance_requirements
ADD compliance_clause_id INT NULL;

-- Add foreign key constraint for compliance_clause_id
ALTER TABLE meeting_compliance_requirements
ADD CONSTRAINT FK_meeting_compliance_requirements_compliance_clauses
FOREIGN KEY (compliance_clause_id)
REFERENCES compliance_clauses(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop foreign key constraint for compliance_requirement_id
ALTER TABLE meeting_compliance_requirements
DROP CONSTRAINT FK_meeting_compliance_requirements_compliance_clauses;



-- Drop compliance_clause_id column
ALTER TABLE meeting_compliance_requirements
DROP COLUMN compliance_clause_id;

-- +goose StatementEnd
