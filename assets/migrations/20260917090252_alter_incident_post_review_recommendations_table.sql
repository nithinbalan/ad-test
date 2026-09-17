-- +goose Up
-- +goose StatementBegin
ALTER TABLE incident_post_review_recommendations ALTER COLUMN recommendation NVARCHAR(MAX) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE incident_post_review_recommendations ALTER COLUMN recommendation NVARCHAR(500) NOT NULL;
-- +goose StatementEnd
