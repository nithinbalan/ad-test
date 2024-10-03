-- +goose Up
-- +goose StatementBegin

ALTER TABLE audit_plans
ADD published_by_id INT NULL;

ALTER TABLE audit_plans
ADD published_at DATETIME NULL;

-- Adding foreign key for published_by_id
ALTER TABLE audit_plans
ADD CONSTRAINT FK_audit_plans_published_by_id FOREIGN KEY (published_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin

-- Dropping foreign key constraint for published_by_id
ALTER TABLE audit_plans
DROP CONSTRAINT FK_audit_plans_published_by_id;

-- Dropping the new fields
ALTER TABLE audit_plans
DROP COLUMN published_by_id;

ALTER TABLE audit_plans
DROP COLUMN published_at;

-- +goose StatementEnd
