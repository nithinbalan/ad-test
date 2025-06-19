-- +goose Up
-- +goose StatementBegin
CREATE TABLE countries
(
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(250) NOT NULL,
    title_ar NVARCHAR(250) NOT NULL,
    code VARCHAR(10) NOT NULL,   
    phone_code VARCHAR(10) NULL,   
    status VARCHAR(10) NOT NULL DEFAULT 'active',

    CONSTRAINT CHK_countries_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE countries;
-- +goose StatementEnd
