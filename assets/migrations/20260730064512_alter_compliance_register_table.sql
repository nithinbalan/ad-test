-- +goose Up
-- +goose StatementBegin

ALTER TABLE compliance_register
ADD
    is_cyber BIT NOT NULL
        CONSTRAINT DF_compliance_register_is_cyber DEFAULT (0);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE compliance_register
DROP CONSTRAINT DF_compliance_register_is_cyber;

ALTER TABLE compliance_register
DROP COLUMN is_cyber;

-- +goose StatementEnd