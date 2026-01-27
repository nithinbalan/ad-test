-- +goose Up
-- +goose StatementBegin
CREATE TABLE issuing_authorities (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en VARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    description_en VARCHAR(255) NOT NULL,
    description_ar NVARCHAR(255) NOT NULL,
    status VARCHAR(10) NOT NULL,
    CONSTRAINT CHK_issuing_authorities_status
        CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.issuing_authorities;
-- +goose StatementEnd
