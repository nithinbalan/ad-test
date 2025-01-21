-- +goose Up
-- +goose StatementBegin
ALTER TABLE controls
ALTER COLUMN control_category_id INT NULL;

ALTER TABLE controls
ALTER COLUMN title_en NVARCHAR(MAX) NOT NULL;

ALTER TABLE controls
ALTER COLUMN title_ar NVARCHAR(MAX) NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE controls
ALTER COLUMN control_category_id INT NOT NULL;

ALTER TABLE controls
ALTER COLUMN title_en NVARCHAR(500) NOT NULL;

ALTER TABLE controls
ALTER COLUMN title_ar NVARCHAR(500) NOT NULL;
-- +goose StatementEnd
