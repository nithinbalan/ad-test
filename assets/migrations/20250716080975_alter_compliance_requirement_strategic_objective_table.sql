-- +goose Up
-- +goose StatementBegin

-- Make compliance_requirement_id nullable
ALTER TABLE compliance_requirement_strategic_objectives
ALTER COLUMN compliance_requirement_id INT NULL;

-- Add new NOT NULL column (no default — works only if table has no existing rows)
ALTER TABLE compliance_requirement_strategic_objectives
ADD compliance_clause_id INT NULL;

ALTER TABLE compliance_requirement_strategic_objectives
ADD CONSTRAINT FK_compliance_requirement_strategic_objective_compliance_clauses
FOREIGN KEY (compliance_clause_id)
REFERENCES compliance_clauses(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE compliance_requirement_strategic_objectives
DROP CONSTRAINT FK_compliance_requirement_strategic_objective_compliance_clauses;

-- Revert compliance_requirement_id to NOT NULL
ALTER TABLE compliance_requirement_strategic_objectives
ALTER COLUMN compliance_requirement_id INT NOT NULL;

-- Drop the compliance_clause_id column
ALTER TABLE compliance_requirement_strategic_objectives
DROP COLUMN compliance_clause_id;

-- +goose StatementEnd
