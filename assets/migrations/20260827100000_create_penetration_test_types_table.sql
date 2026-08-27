-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_types (
    id INT PRIMARY KEY IDENTITY(1, 1),
    code NVARCHAR(20) NOT NULL,
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT UQ_penetration_test_types_code UNIQUE (code),
    CONSTRAINT CHK_penetration_test_types_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS penetration_test_types;
-- +goose StatementEnd
