-- +goose Up
-- +goose StatementBegin
CREATE TABLE kpi_library (
    id INT PRIMARY KEY IDENTITY(1,1),
    category NVARCHAR(255) NOT NULL,
    title NVARCHAR(500) NOT NULL,
    description NVARCHAR(MAX) NULL,
    unit NVARCHAR(50) NOT NULL,
    frequency NVARCHAR(50) NOT NULL
);
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS kpi_library;
-- +goose StatementEnd