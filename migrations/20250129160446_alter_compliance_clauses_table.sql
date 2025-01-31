-- +goose Up
-- +goose StatementBegin
-- Rename 'title' to 'title_en'
EXEC sp_rename 'compliance_clauses.title', 'title_en', 'COLUMN';

-- Add each column separately
ALTER TABLE compliance_clauses ADD compliance_clause_id INT NULL;
ALTER TABLE compliance_clauses ADD title_ar NVARCHAR(MAX) NULL;
ALTER TABLE compliance_clauses ADD description_en NVARCHAR(MAX) NULL;
ALTER TABLE compliance_clauses ADD description_ar NVARCHAR(MAX) NULL;
ALTER TABLE compliance_clauses ADD scope_en NVARCHAR(MAX) NULL;
ALTER TABLE compliance_clauses ADD scope_ar NVARCHAR(MAX) NULL;
ALTER TABLE compliance_clauses ADD objective_en NVARCHAR(MAX) NULL;
ALTER TABLE compliance_clauses ADD objective_ar NVARCHAR(MAX) NULL;
ALTER TABLE compliance_clauses ADD created_by_id INT NULL;
ALTER TABLE compliance_clauses ADD created_at DATETIME NULL;
ALTER TABLE compliance_clauses ADD updated_by_id INT NULL;
ALTER TABLE compliance_clauses ADD updated_at DATETIME NULL;

-- Add foreign key constraints separately
ALTER TABLE compliance_clauses ADD CONSTRAINT FK_compliance_clauses_self FOREIGN KEY (compliance_clause_id) REFERENCES compliance_clauses(id);
ALTER TABLE compliance_clauses ADD CONSTRAINT FK_compliance_clauses_user_1 FOREIGN KEY (created_by_id) REFERENCES users(id);
ALTER TABLE compliance_clauses ADD CONSTRAINT FK_compliance_clauses_user_2 FOREIGN KEY (updated_by_id) REFERENCES users(id);
-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin
-- Drop foreign key constraints first
ALTER TABLE compliance_clauses DROP CONSTRAINT FK_compliance_clauses_self;
ALTER TABLE compliance_clauses DROP CONSTRAINT FK_compliance_clauses_user_1;
ALTER TABLE compliance_clauses DROP CONSTRAINT FK_compliance_clauses_user_2;

-- Drop each column separately
ALTER TABLE compliance_clauses DROP COLUMN compliance_clause_id;
ALTER TABLE compliance_clauses DROP COLUMN title_ar;
ALTER TABLE compliance_clauses DROP COLUMN description_en;
ALTER TABLE compliance_clauses DROP COLUMN description_ar;
ALTER TABLE compliance_clauses DROP COLUMN scope_en;
ALTER TABLE compliance_clauses DROP COLUMN scope_ar;
ALTER TABLE compliance_clauses DROP COLUMN objective_en;
ALTER TABLE compliance_clauses DROP COLUMN objective_ar;
ALTER TABLE compliance_clauses DROP COLUMN created_by_id;
ALTER TABLE compliance_clauses DROP COLUMN created_at;
ALTER TABLE compliance_clauses DROP COLUMN updated_by_id;
ALTER TABLE compliance_clauses DROP COLUMN updated_at;

-- Rename 'title_en' back to 'title'
EXEC sp_rename 'compliance_clauses.title_en', 'title', 'COLUMN';
-- +goose StatementEnd


