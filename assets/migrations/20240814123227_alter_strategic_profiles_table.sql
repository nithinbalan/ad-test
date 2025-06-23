-- +goose Up
-- +goose StatementBegin
ALTER TABLE strategic_profiles
ADD profile_image NVARCHAR(255) NULL,
    profile_image_token NVARCHAR(255) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE strategic_profiles
DROP COLUMN profile_image,
DROP COLUMN profile_image_token;
-- +goose StatementEnd
