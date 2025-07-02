-- +goose Up
-- +goose StatementBegin
CREATE TABLE swot_pestel_categories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL,
        CONSTRAINT CHK_swot_pestel_categories_type CHECK (type IN ('strength', 'weakness','opportunity','threat')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_swot_pestel_categories_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.swot_pestel_categories;
-- +goose StatementEnd

