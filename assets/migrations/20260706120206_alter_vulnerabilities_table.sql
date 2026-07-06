-- +goose Up
-- +goose StatementBegin

ALTER TABLE vulnerabilities
ADD 
    external_id VARCHAR(100) NULL,
    external_metadata NVARCHAR(MAX) NULL;

-- Create JSON check constraint using dynamic SQL
EXEC('
ALTER TABLE vulnerabilities
ADD CONSTRAINT CK_vulnerabilities_external_metadata_json
CHECK (external_metadata IS NULL OR ISJSON(external_metadata) = 1);
');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
EXEC('
ALTER TABLE vulnerabilities
DROP CONSTRAINT CK_vulnerabilities_external_metadata_json;
');

ALTER TABLE vulnerabilities
DROP COLUMN external_id, external_metadata;
-- +goose StatementEnd
