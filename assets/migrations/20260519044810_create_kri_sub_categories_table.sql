-- +goose Up
-- +goose StatementBegin

CREATE TABLE kri_sub_categories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    kri_category_id INT NOT NULL,
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',

    CONSTRAINT CHK_kri_sub_categories_status 
        CHECK (status IN ('active', 'inactive')),

    CONSTRAINT FK_kri_sub_categories_kri_categories 
        FOREIGN KEY (kri_category_id) 
        REFERENCES kri_categories(id) 
        ON DELETE NO ACTION 
        ON UPDATE NO ACTION
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DROP TABLE dbo.kri_sub_categories;

-- +goose StatementEnd