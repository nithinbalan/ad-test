-- +goose Up
-- +goose StatementBegin

ALTER TABLE consultation_request_comments
ADD is_edited BIT NOT NULL CONSTRAINT df_consultation_request_comments_is_edited DEFAULT 0;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE consultation_request_comments
DROP CONSTRAINT df_consultation_request_comments_is_edited;


ALTER TABLE consultation_request_comments
DROP COLUMN is_edited;

-- +goose StatementEnd
