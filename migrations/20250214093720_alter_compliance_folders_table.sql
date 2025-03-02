-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_folders ADD created_at DATETIME NULL;
ALTER TABLE compliance_folders ADD updated_at DATETIME NULL;
ALTER TABLE compliance_folders ADD created_by_id INT NULL;
ALTER TABLE compliance_folders ADD updated_by_id INT NULL;

-- Add foreign key constraints separately
ALTER TABLE compliance_folders
ADD CONSTRAINT FK_compliance_folders_user_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

ALTER TABLE compliance_folders
ADD CONSTRAINT FK_compliance_folders_user_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Drop foreign key constraints first
ALTER TABLE compliance_folders DROP CONSTRAINT FK_compliance_folders_user_1;
ALTER TABLE compliance_folders DROP CONSTRAINT FK_compliance_folders_user_2;

-- Drop columns separately
ALTER TABLE compliance_folders DROP COLUMN created_at;
ALTER TABLE compliance_folders DROP COLUMN updated_at;
ALTER TABLE compliance_folders DROP COLUMN created_by_id;
ALTER TABLE compliance_folders DROP COLUMN updated_by_id;
-- +goose StatementEnd
