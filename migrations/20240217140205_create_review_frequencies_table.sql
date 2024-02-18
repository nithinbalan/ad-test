-- +goose Up
-- +goose StatementBegin
CREATE TABLE review_frequencies (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type NVARCHAR(50) NOT NULL,
    CONSTRAINT CHK_ReviewFrequencyType CHECK (
        type IN (
            'daily',
            'weekly',
            'monthly',
            'quarterly',
            'half-yearly',
            'annually'
        )
    ) -- Allowed values
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.review_frequencies;

-- +goose StatementEnd