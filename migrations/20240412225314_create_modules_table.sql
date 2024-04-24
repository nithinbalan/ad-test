-- +goose Up
-- +goose StatementBegin
CREATE TABLE modules (
    id INT PRIMARY KEY IDENTITY(1,1),
    identifier VARCHAR(45) NOT NULL,
    title_en VARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    module_id INT,
    [order] INT,
    CONSTRAINT FK_modules_module_id FOREIGN KEY (module_id) REFERENCES modules(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS modules;
-- +goose StatementEnd
