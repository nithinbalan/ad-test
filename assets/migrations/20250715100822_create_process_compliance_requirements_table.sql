-- +goose Up
-- +goose StatementBegin
CREATE TABLE process_compliance_requirements (
  id INT PRIMARY KEY IDENTITY(1,1),
  process_id INT NOT NULL,
  compliance_register_id INT NOT NULL,
  compliance_requirement_id INT NULL DEFAULT NULL,
  compliance_clause_id INT NOT NULL,
  created_by_id INT,
  created_at DATETIME NOT NULL,
  CONSTRAINT FKprocess_compliance_requirement_process FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT FK_process_compliance_requirements_compliance_register FOREIGN KEY (compliance_register_id) REFERENCES compliance_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT FK_process_compliance_requirement_compliance_clauses FOREIGN KEY (compliance_clause_id) REFERENCES compliance_clauses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_process_compliance_requirement_compliance_requirement FOREIGN KEY (compliance_requirement_id) REFERENCES compliance_requirements(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT FK_process_compliance_requirement_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.process_compliance_requirements;
-- +goose StatementEnd
