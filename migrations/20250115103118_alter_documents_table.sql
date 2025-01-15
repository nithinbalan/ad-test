-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents
ADD next_review_date DATE NULL,
    review_frequency_id INT NULL,
    CONSTRAINT FK_documents_review_frequencies FOREIGN KEY (review_frequency_id) REFERENCES review_frequencies(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents DROP CONSTRAINT FK_documents_review_frequencies;
ALTER TABLE documents DROP COLUMN next_review_date, DROP COLUMN review_frequency_id;
-- +goose StatementEnd
