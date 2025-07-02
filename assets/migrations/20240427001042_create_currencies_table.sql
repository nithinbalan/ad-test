-- +goose Up
-- +goose StatementBegin
CREATE TABLE currencies (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title VARCHAR(45) NOT NULL,
    code VARCHAR(10) NOT NULL,
    icon VARCHAR(45) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS currencies;
-- +goose StatementEnd
