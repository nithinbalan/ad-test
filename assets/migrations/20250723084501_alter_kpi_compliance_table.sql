-- +goose Up
-- +goose StatementBegin

-- Add nullable compliance_requirement_id column
ALTER TABLE kpi_compliances
ADD compliance_requirement_id INT NULL;

-- Add foreign key constraint for compliance_requirement_id
ALTER TABLE kpi_compliances
ADD CONSTRAINT FK_kpi_compliances_compliance_requirements
FOREIGN KEY (compliance_requirement_id)
REFERENCES compliance_requirements(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- Add nullable compliance_clause_id column
ALTER TABLE kpi_compliances
ADD compliance_clause_id INT NULL;

-- Add foreign key constraint for compliance_clause_id
ALTER TABLE kpi_compliances
ADD CONSTRAINT FK_kpi_compliances_compliance_clauses
FOREIGN KEY (compliance_clause_id)
REFERENCES compliance_clauses(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Drop foreign key constraint for compliance_requirement_id
ALTER TABLE kpi_compliances
DROP CONSTRAINT FK_kpi_compliances_compliance_requirements;

-- Drop compliance_requirement_id column
ALTER TABLE kpi_compliances
DROP COLUMN compliance_requirement_id;

-- Drop foreign key constraint for compliance_clause_id
ALTER TABLE kpi_compliances
DROP CONSTRAINT FK_kpi_compliances_compliance_clauses;

-- Drop compliance_clause_id column
ALTER TABLE kpi_compliances
DROP COLUMN compliance_clause_id;

-- +goose StatementEnd
