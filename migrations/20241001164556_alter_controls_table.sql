-- +goose Up
-- +goose StatementBegin
ALTER TABLE controls
ADD 
    is_applicable BIT NOT NULL DEFAULT 0,
    justification NVARCHAR(MAX),
    applicability_change_date DATETIME;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE controls
DROP COLUMN is_applicable,
DROP COLUMN justification,
DROP COLUMN applicability_change_date;
-- +goose StatementEnd
