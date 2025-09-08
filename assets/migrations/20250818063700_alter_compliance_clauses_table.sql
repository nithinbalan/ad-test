-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_clauses
ADD compliance_status_id INT NULL,
    score DECIMAL(10,2) NULL;
 
ALTER TABLE compliance_clauses
ADD CONSTRAINT FK_compliance_clauses_compliance_statuses
FOREIGN KEY (compliance_status_id) REFERENCES compliance_statuses(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_clauses
DROP CONSTRAINT FK_compliance_clauses_compliance_statuses;
 
ALTER TABLE compliance_clauses
DROP COLUMN score;
ALTER TABLE compliance_clauses
    DROP COLUMN compliance_status_id;
-- +goose StatementEnd