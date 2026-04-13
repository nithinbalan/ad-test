-- +goose Up
-- +goose StatementBegin

ALTER TABLE internal_committee_review_comments
ADD 
    is_edited BIT NOT NULL 
        CONSTRAINT DF_internal_committee_review_comments_is_edited DEFAULT 0;

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

-- Drop DEFAULT constraint
ALTER TABLE internal_committee_review_comments
DROP CONSTRAINT DF_internal_committee_review_comments_is_edited;

-- Then drop column
ALTER TABLE internal_committee_review_comments
DROP COLUMN is_edited;

-- +goose StatementEnd