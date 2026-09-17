-- +goose Up
-- +goose StatementBegin
ALTER TABLE incident_post_review_lessons ALTER COLUMN lesson NVARCHAR(MAX) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE incident_post_review_lessons ALTER COLUMN lesson NVARCHAR(500) NOT NULL;
-- +goose StatementEnd
