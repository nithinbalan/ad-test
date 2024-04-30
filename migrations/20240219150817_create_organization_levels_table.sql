-- +goose Up
-- +goose StatementBegin
CREATE TABLE organization_levels (
    id INT PRIMARY KEY IDENTITY,
    level TINYINT NOT NULL CHECK (level BETWEEN 1 AND 10),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    is_primary TINYINT NOT NULL DEFAULT 0 CHECK (is_primary IN (0, 1))
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.organization_levels;

-- +goose StatementEnd