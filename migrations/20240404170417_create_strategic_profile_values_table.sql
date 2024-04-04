-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategic_profile_values (
    id INT IDENTITY(1,1) PRIMARY KEY,
    strategic_profile_id INT NOT NULL,
    title_en VARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en VARCHAR(500),
    description_ar VARCHAR(500),
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    CONSTRAINT FK_strategic_profile_values_strategic_profile_id FOREIGN KEY (strategic_profile_id) REFERENCES strategic_profiles(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.strategic_profile_values;
-- +goose StatementEnd
