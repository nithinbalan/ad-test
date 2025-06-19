-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents
ADD review_status_id INT NULL;

ALTER TABLE documents
ADD CONSTRAINT FK_documents_review_statuses FOREIGN KEY (review_status_id) REFERENCES review_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents DROP CONSTRAINT FK_documents_review_statuses;

ALTER TABLE documents DROP COLUMN review_status_id;
-- +goose StatementEnd
