-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_change_requests
ADD document_review_id INT NULL;

ALTER TABLE document_change_requests
ADD CONSTRAINT FK_document_cr_document_review FOREIGN KEY (document_review_id) REFERENCES document_reviews(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE document_change_requests
DROP CONSTRAINT FK_document_cr_document_review;

ALTER TABLE document_change_requests
DROP COLUMN document_review_id;
-- +goose StatementEnd
