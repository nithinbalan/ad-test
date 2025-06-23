-- +goose Up
-- +goose StatementBegin
-- Rename 'title' to 'title_en'
EXEC sp_rename 'compliance_requirements.title', 'title_en', 'COLUMN';

-- Add 'title_ar' column
ALTER TABLE compliance_requirements
    ADD title_ar NVARCHAR(MAX) NULL;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
-- Drop 'title_ar' column
ALTER TABLE compliance_requirements
    DROP COLUMN title_ar;

-- Rename 'title_en' back to 'title'
EXEC sp_rename 'compliance_requirements.title_en', 'title', 'COLUMN';
-- +goose StatementEnd
