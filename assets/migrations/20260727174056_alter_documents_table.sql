-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents
ADD confidentiality_classification_id INT NULL;

ALTER TABLE documents ADD CONSTRAINT FK_documents_confidentiality_classifications
    FOREIGN KEY (confidentiality_classification_id) REFERENCES confidentiality_classifications(id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents
DROP CONSTRAINT FK_documents_confidentiality_classifications;

ALTER TABLE documents
DROP COLUMN confidentiality_classification_id;
-- +goose StatementEnd
