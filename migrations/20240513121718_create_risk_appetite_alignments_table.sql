-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_appetite_alignments (
  id INT PRIMARY KEY IDENTITY(1,1),
  risk_evaluation_id INT NOT NULL,
  risk_appetite_statement_id INT NOT NULL,
  risk_appetite_tolerance_level_id INT NOT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NULL,
  created_by INT NOT NULL,
  updated_by INT NULL,
  CONSTRAINT FK_risk_appetite_alignments_risk_evaluation FOREIGN KEY (risk_evaluation_id) REFERENCES risk_evaluations(id),
  CONSTRAINT FK_risk_appetite_alignments_risk_appetite_statements FOREIGN KEY (risk_appetite_statement_id) REFERENCES risk_appetite_statements(id),
  CONSTRAINT FK_risk_appetite_alignments_risk_appetite_tolerance_levels FOREIGN KEY (risk_appetite_tolerance_level_id) REFERENCES risk_appetite_tolerance_levels(id),
  CONSTRAINT FK_risk_appetite_alignments_users_1 FOREIGN KEY (created_by) REFERENCES users(id),
  CONSTRAINT FK_risk_appetite_alignments_users_2 FOREIGN KEY (updated_by) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS risk_appetite_alignments;
-- +goose StatementEnd
