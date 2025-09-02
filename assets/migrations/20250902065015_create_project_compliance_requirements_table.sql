-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_compliance_requirements (
  id INT PRIMARY KEY IDENTITY(1,1),
  project_id INT NOT NULL,
  compliance_register_id INT NOT NULL,
  compliance_requirement_id INT NULL DEFAULT NULL,
  compliance_clause_id INT NOT NULL,
  created_by_id INT,
  created_at DATETIME NOT NULL,
  CONSTRAINT FK_project_compliance_requirement_projects FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT FK_project_compliance_requirement_compliance_register FOREIGN KEY (compliance_register_id) REFERENCES compliance_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT FK_project_compliance_requirement_compliance_clauses FOREIGN KEY (compliance_clause_id) REFERENCES compliance_clauses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK__project_compliance_requirement_compliance_requirement FOREIGN KEY (compliance_requirement_id) REFERENCES compliance_requirements(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT FK_project_compliance_requirement_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.project_compliance_requirements;
-- +goose StatementEnd
