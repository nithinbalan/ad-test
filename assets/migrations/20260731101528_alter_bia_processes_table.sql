-- +goose Up
-- +goose StatementBegin

ALTER TABLE bia_processes
ADD
    is_cyber BIT NOT NULL
        CONSTRAINT DF_bia_processes_is_cyber DEFAULT (0);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE bia_processes
DROP CONSTRAINT DF_bia_processes_is_cyber;

ALTER TABLE bia_processes
DROP COLUMN is_cyber;

-- +goose StatementEnd