-- +goose Up
-- +goose StatementBegin
CREATE TABLE time_units (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_time_units_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
