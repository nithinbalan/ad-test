-- +goose Up
-- +goose StatementBegin

ALTER TABLE strategic_profile_values
ALTER COLUMN title_en NVARCHAR(MAX) NOT NULL;

ALTER TABLE strategic_profile_values
ALTER COLUMN title_ar NVARCHAR(MAX) NOT NULL;

ALTER TABLE strategic_profile_values
ALTER COLUMN description_en NVARCHAR(MAX) NULL;

ALTER TABLE strategic_profile_values
ALTER COLUMN description_ar NVARCHAR(MAX) NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE strategic_profile_values
ALTER COLUMN title_en NVARCHAR(255) NOT NULL;

ALTER TABLE strategic_profile_values
ALTER COLUMN title_ar NVARCHAR(255) NOT NULL;

ALTER TABLE strategic_profile_values
ALTER COLUMN description_en NVARCHAR(500) NULL;

ALTER TABLE strategic_profile_values
ALTER COLUMN description_ar NVARCHAR(500) NULL;

-- +goose StatementEnd