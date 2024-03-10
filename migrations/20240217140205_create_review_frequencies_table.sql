-- +goose Up
-- +goose StatementBegin
CREATE TABLE review_frequencies (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en VARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(45) NOT NULL,
        CONSTRAINT CHK_review_frequencies_type CHECK (type IN ('daily', 'weekly', 'monthly', 'quarterly', 'half-yearly', 'annually')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_review_frequencies_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.review_frequencies;
-- +goose StatementEnd