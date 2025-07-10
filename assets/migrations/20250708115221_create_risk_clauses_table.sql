-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_clauses (
  id INT PRIMARY KEY IDENTITY(1,1),
  risk_register_id INT NOT NULL,
  compliance_register_id INT NOT NULL,
  compliance_clause_id INT NOT NULL,
  created_by_id INT,
  created_at DATETIME NOT NULL,
  CONSTRAINT FK_risk_clause_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT FK_risk_clause_compliance_register FOREIGN KEY (compliance_register_id) REFERENCES compliance_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT FK_risk_clause_compliance_clauses FOREIGN KEY (compliance_clause_id) REFERENCES compliance_clauses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_risk_clauses_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.risk_clauses;
-- +goose StatementEnd
