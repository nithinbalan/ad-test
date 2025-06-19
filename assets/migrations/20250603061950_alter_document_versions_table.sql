-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_versions
ADD parsed_content NVARCHAR(MAX) NULL,
    sample_questions NVARCHAR(MAX) NULL,
    summary NVARCHAR(MAX) NULL,
    key_points NVARCHAR(MAX) NULL,
    faq NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE document_versions
DROP COLUMN sample_questions,
            parsed_content,
            summary,
            key_points,
            faq;
-- +goose StatementEnd 