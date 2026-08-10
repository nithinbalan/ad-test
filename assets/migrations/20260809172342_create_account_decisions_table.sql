-- +goose Up
-- +goose StatementBegin
CREATE TABLE account_decisions (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(20) NOT NULL,
    status VARCHAR(10) NOT NULL,

    CONSTRAINT UQ_account_decisions_type UNIQUE (type),
    CONSTRAINT CHK_account_decisions_type
        CHECK (type IN ('keep', 'revoke', 'modify', 'investigate')),
    CONSTRAINT CHK_account_decisions_status
        CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE account_decisions;
-- +goose StatementEnd
