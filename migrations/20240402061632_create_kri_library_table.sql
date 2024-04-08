-- +goose Up
-- +goose StatementBegin
CREATE TABLE kri_library (
    id INT PRIMARY KEY IDENTITY(1, 1),
    kri_category_id INT NOT NULL,
    review_frequency_id INT NULL,
    title NVARCHAR(255) NOT NULL,
    description NVARCHAR(MAX) NULL, 
    calculation_type VARCHAR(10) NOT NULL DEFAULT 'automatic',
    CONSTRAINT CHK_kri_library_calculation_type CHECK (calculation_type IN ('automatic', 'manual')),
    source VARCHAR(45) NOT NULL,
    CONSTRAINT FK_kri_category_kri_library FOREIGN KEY (kri_category_id) REFERENCES kri_categories(id),
    CONSTRAINT FK_review_frequency_kri_library FOREIGN KEY (review_frequency_id) REFERENCES review_frequencies(id),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_kri_library_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.kri_library;
-- +goose StatementEnd
