-- +goose Up
-- +goose StatementBegin
CREATE TABLE feedback_options (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL DEFAULT 'poor',
        CONSTRAINT CHK_feedback_options_type CHECK (type IN ('poor', 'fair', 'good', 'excellent')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_feedback_options_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.feedback_options;
-- +goose StatementEnd

