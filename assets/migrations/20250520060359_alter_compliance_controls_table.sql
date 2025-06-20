-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_controls
ADD compliance_clause_id INT NULL;

ALTER TABLE compliance_controls
ADD CONSTRAINT FK_compliance_controls_compliance_clauses
FOREIGN KEY (compliance_clause_id) REFERENCES compliance_clauses(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_controls
DROP CONSTRAINT FK_compliance_controls_compliance_clauses;

ALTER TABLE compliance_controls
DROP COLUMN compliance_clause_id;
-- +goose StatementEnd
