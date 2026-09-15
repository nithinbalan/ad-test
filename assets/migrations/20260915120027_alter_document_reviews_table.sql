-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_reviews
ADD target_review_date DATE NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE document_reviews
DROP COLUMN target_review_date;
-- +goose StatementEnd


