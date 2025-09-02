-- +goose Up
-- +goose StatementBegin
CREATE TABLE kri_compliance_requirements (
  id INT PRIMARY KEY IDENTITY(1,1),
  kri_register_id INT NOT NULL,
  compliance_register_id INT NOT NULL,
  compliance_requirement_id INT NULL DEFAULT NULL,
  compliance_clause_id INT NOT NULL,
  created_by_id INT,
  created_at DATETIME NOT NULL,
  CONSTRAINT FK_kri_compliance_requirement_kri_register FOREIGN KEY (kri_register_id) REFERENCES kri_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT FK_kri_compliance_requirement_compliance_register FOREIGN KEY (compliance_register_id) REFERENCES compliance_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT FK_kri_compliance_requirement_compliance_clauses FOREIGN KEY (compliance_clause_id) REFERENCES compliance_clauses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK__kri_compliance_requirement_compliance_requirement FOREIGN KEY (compliance_requirement_id) REFERENCES compliance_requirements(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
  CONSTRAINT FK_kri_compliance_requirement_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.kri_compliance_requirements;
-- +goose StatementEnd
