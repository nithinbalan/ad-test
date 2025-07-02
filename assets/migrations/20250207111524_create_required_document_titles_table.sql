-- +goose Up
-- +goose StatementBegin
CREATE TABLE required_document_titles (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    status VARCHAR(10) NOT NULL,
        CONSTRAINT CHK_required_document_titles_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.required_document_titles;
-- +goose StatementEnd
