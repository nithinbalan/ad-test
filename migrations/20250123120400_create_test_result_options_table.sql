-- +goose Up
-- +goose StatementBegin
CREATE TABLE test_result_options (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(25) NOT NULL DEFAULT 'successful',
        CONSTRAINT CHK_test_result_options_type CHECK (type IN ('successful', 'partially-successful', 'failure')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_test_result_options_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.test_result_options;
-- +goose StatementEnd

