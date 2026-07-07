-- +goose Up
-- +goose StatementBegin
ALTER TABLE severity_levels
DROP CONSTRAINT CHK_severity_levels_type;

ALTER TABLE severity_levels
ADD CONSTRAINT CHK_severity_levels_type
CHECK (type IN ('critical', 'high', 'medium', 'low'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE severity_levels
DROP CONSTRAINT CHK_severity_levels_type;

ALTER TABLE severity_levels
ADD CONSTRAINT CHK_severity_levels_type
CHECK (type IN ('high', 'medium', 'low'));
-- +goose StatementEnd
