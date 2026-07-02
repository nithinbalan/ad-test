-- +goose Up
-- +goose StatementBegin

ALTER TABLE violation_reports
ADD is_external BIT NOT NULL
    CONSTRAINT DF_violation_reports_is_external DEFAULT (0);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE violation_reports
DROP CONSTRAINT DF_violation_reports_is_external;

ALTER TABLE violation_reports
DROP COLUMN is_external;

-- +goose StatementEnd