-- +goose Up
-- +goose StatementBegin
-- Add the vulnerability_id column with a foreign key to vulnerabilities table
ALTER TABLE threats
ADD vulnerability_id INT NULL;

-- Add the foreign key constraint for vulnerability_id to vulnerabilities table
ALTER TABLE threats
ADD CONSTRAINT FK_threats_vulnerability_id
FOREIGN KEY (vulnerability_id) REFERENCES vulnerabilities(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove the foreign key constraint for vulnerability_id
ALTER TABLE threats
DROP CONSTRAINT FK_threats_vulnerability_id;

-- Drop the vulnerability_id column
ALTER TABLE threats
DROP COLUMN vulnerability_id;
-- +goose StatementEnd
