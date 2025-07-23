-- +goose Up
-- +goose StatementBegin

-- Make compliance_requirement_id nullable
ALTER TABLE kpi_compliances
ADD compliance_requirement_id INT NULL;

ALTER TABLE kpi_compliances
ADD CONSTRAINT FK_kpi_compliances_compliance_requirements
FOREIGN KEY (compliance_requirement_id)
REFERENCES compliance_requirements(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;

-- Add new NOT NULL column
ALTER TABLE kpi_compliances
ADD compliance_clause_id INT NULL;

ALTER TABLE kpi_compliances
ADD CONSTRAINT FK_kpi_compliances_compliance_clauses
FOREIGN KEY (compliance_clause_id)
REFERENCES compliance_clauses(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE kpi_compliances
DROP CONSTRAINT  FK_kpi_compliances_compliance_register;

-- Drop the compliance_register_id column
ALTER TABLE kpi_compliances
DROP COLUMN compliance_register_id;

ALTER TABLE kpi_compliances
DROP CONSTRAINT FK_kpi_compliances_compliance_clauses;

-- Drop the compliance_clause_id column
ALTER TABLE kpi_compliances
DROP COLUMN compliance_clause_id;

-- +goose StatementEnd