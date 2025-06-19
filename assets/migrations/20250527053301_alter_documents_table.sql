-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents
ADD 
    document_content NVARCHAR(MAX) NULL,
    sample_questions NVARCHAR(MAX) NOT NULL DEFAULT N'{}';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE documents
DROP COLUMN sample_questions,
DROP COLUMN document_content;
-- +goose StatementEnd
