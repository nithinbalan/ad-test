-- +goose Up
-- +goose StatementBegin
CREATE TABLE ai_modules (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NULL,
    [order] INT DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS ai_modules;
-- +goose StatementEnd
