-- +goose Up
-- +goose StatementBegin
CREATE TABLE timezones (
    id INT PRIMARY KEY IDENTITY(1,1),
    name VARCHAR(255),
    offset VARCHAR(100),
    display_name VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.timezones;
-- +goose StatementEnd
