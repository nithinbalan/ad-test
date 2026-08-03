-- +goose Up
-- +goose StatementBegin

ALTER TABLE business_functions
ADD
    is_cyber BIT NOT NULL
        CONSTRAINT DF_business_functions_is_cyber DEFAULT (0);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE business_functions
DROP CONSTRAINT DF_business_functions_is_cyber;

ALTER TABLE business_functions
DROP COLUMN is_cyber;

-- +goose StatementEnd