-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_library (
    id INT PRIMARY KEY IDENTITY(1, 1),
    category NVARCHAR(255) NOT NULL,
    title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX) NULL,
    likelihood NVARCHAR(50) NOT NULL,
    impact NVARCHAR(50) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS risk_library;
-- +goose StatementEnd
