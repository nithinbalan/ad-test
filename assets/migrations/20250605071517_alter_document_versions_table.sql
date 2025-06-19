-- +goose Up
-- +goose StatementBegin
ALTER TABLE document_versions
DROP COLUMN sample_questions,
            summary,
            key_points,
            faq;

ALTER TABLE document_versions
ADD sample_questions_en NVARCHAR(MAX) NULL,
    summary_en NVARCHAR(MAX) NULL,
    key_points_en NVARCHAR(MAX) NULL,
    faq_en NVARCHAR(MAX) NULL,
    sample_questions_ar NVARCHAR(MAX) NULL,
    summary_ar NVARCHAR(MAX) NULL,
    key_points_ar NVARCHAR(MAX) NULL,
    faq_ar NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE document_versions
DROP COLUMN sample_questions_en,
            summary_en,
            key_points_en,
            faq_en,
            sample_questions_ar,
            summary_ar,
            key_points_ar,
            faq_ar;

ALTER TABLE document_versions
ADD sample_questions NVARCHAR(MAX) NULL,
    summary NVARCHAR(MAX) NULL,
    key_points NVARCHAR(MAX) NULL,
    faq NVARCHAR(MAX) NULL;
-- +goose StatementEnd
