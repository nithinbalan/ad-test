-- +goose Up
-- +goose StatementBegin
CREATE TABLE lifecycle_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en VARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(45) NOT NULL,
    type VARCHAR(45) NOT NULL DEFAULT 'not_implemented',
        CONSTRAINT CHK_lifecycle_statuses_type CHECK (type IN ('active', 'expired', 'archived')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_lifecycle_statuses_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.lifecycle_statuses;
-- +goose StatementEnd
