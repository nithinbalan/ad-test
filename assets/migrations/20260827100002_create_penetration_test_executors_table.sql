-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_executors (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT UQ_penetration_test_executors_title_en UNIQUE (title_en),
    CONSTRAINT CHK_penetration_test_executors_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS penetration_test_executors;
-- +goose StatementEnd
