-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_root_cause_analysis (
  id INT PRIMARY KEY IDENTITY(1,1),
  risk_register_id INT NOT NULL,
  why INT NOT NULL,
  question NVARCHAR(500) NOT NULL,
  answer NVARCHAR(1000) NULL,
  created_by_id INT NOT NULL,
  updated_by_id INT NULL,
  created_at DATETIME NOT NULL,
  updated_at DATETIME NULL,
  CONSTRAINT fk_risk_root_cause_analysis_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id),
  CONSTRAINT fk_risk_root_cause_analysis_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id),
  CONSTRAINT fk_risk_root_cause_analysis_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS risk_root_cause_analysis;
-- +goose StatementEnd
