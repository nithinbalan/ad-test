-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategy_profile_settings (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    is_enable TINYINT NOT NULL DEFAULT 1
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.strategy_profile_settings;
-- +goose StatementEnd
