-- +goose Up
-- +goose StatementBegin
ALTER TABLE control_types
ADD type VARCHAR(20) NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE control_types
ADD CONSTRAINT CHK_control_types_type
CHECK (type IN ('preventive', 'corrective', 'detective', 'deterrent'));
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE control_types
DROP CONSTRAINT CHK_control_types_type;

ALTER TABLE control_types
DROP COLUMN type;
-- +goose StatementEnd
