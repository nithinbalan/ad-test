-- +goose Up
-- +goose StatementBegin
CREATE TABLE control_library (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(255) ,
    title_ar NVARCHAR(255) ,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE control_library;
-- +goose StatementEnd
