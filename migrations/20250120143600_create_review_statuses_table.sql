-- +goose Up
-- +goose StatementBegin
CREATE TABLE review_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(20) NOT NULL DEFAULT 'review-pending',
        CONSTRAINT CHK_review_statuses_type CHECK (type IN ('reviewed', 'review-pending')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_review_statuses_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.review_statuses;
-- +goose StatementEnd
