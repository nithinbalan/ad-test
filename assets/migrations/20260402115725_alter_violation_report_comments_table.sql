-- +goose Up
-- +goose StatementBegin

ALTER TABLE violation_report_comments
ADD is_edited BIT NOT NULL CONSTRAINT df_violation_report_comments_is_edited DEFAULT 0;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE violation_report_comments
DROP CONSTRAINT df_violation_report_comments_is_edited;


ALTER TABLE violation_report_comments
DROP COLUMN is_edited;

-- +goose StatementEnd
