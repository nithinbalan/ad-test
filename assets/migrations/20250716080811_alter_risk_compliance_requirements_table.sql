-- +goose Up
-- +goose StatementBegin

-- Make compliance_requirement_id nullable
ALTER TABLE risk_compliance_requirements
ALTER COLUMN compliance_requirement_id INT NULL;


ALTER TABLE risk_compliance_requirements
ADD compliance_clause_id INT NULL;

ALTER TABLE risk_compliance_requirements
ADD CONSTRAINT FK_risk_compliance_requirement_compliance_clauses
FOREIGN KEY (compliance_clause_id)
REFERENCES compliance_clauses(id)
ON DELETE NO ACTION
ON UPDATE NO ACTION;


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE risk_compliance_requirements
DROP CONSTRAINT FK_risk_compliance_requirement_compliance_clauses;

-- Revert compliance_requirement_id to NOT NULL
ALTER TABLE risk_compliance_requirements
ALTER COLUMN compliance_requirement_id INT NOT NULL;

-- Drop the compliance_clause_id column
ALTER TABLE risk_compliance_requirements
DROP COLUMN compliance_clause_id;

-- +goose StatementEnd
