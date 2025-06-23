-- +goose Up
-- +goose StatementBegin
CREATE TABLE trigger_conditions (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    type VARCHAR(25) NOT NULL,
        CONSTRAINT CHK_trigger_conditions_type CHECK (type IN ('before-due-date','on-due-date','after-due-date')),
    status VARCHAR(10) NOT NULL,
        CONSTRAINT CHK_trigger_conditions_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.trigger_conditions;
-- +goose StatementEnd
