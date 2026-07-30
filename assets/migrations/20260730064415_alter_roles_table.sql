-- +goose Up
-- +goose StatementBegin

ALTER TABLE roles
ADD
    is_cyber BIT NOT NULL
        CONSTRAINT DF_roles_is_cyber DEFAULT (0);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE roles
DROP CONSTRAINT DF_roles_is_cyber;

ALTER TABLE roles
DROP COLUMN is_cyber;

-- +goose StatementEnd