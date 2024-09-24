-- +goose Up
-- +goose StatementBegin
ALTER TABLE audit_schedules
ADD comment NVARCHAR(MAX) NULL,
    percentage INT NOT NULL DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE audit_schedules
DROP COLUMN comment,
    DROP COLUMN percentage;
-- +goose StatementEnd