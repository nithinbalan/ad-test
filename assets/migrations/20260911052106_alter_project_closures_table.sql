-- +goose Up
-- +goose StatementBegin
ALTER TABLE project_closures
ADD is_approved BIT NOT NULL
        CONSTRAINT DF_project_closures_is_approved DEFAULT 0,
    approved_at DATETIME NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE project_closures
DROP CONSTRAINT DF_project_closures_is_approved;

ALTER TABLE project_closures
DROP COLUMN is_approved,
            approved_at;
-- +goose StatementEnd
