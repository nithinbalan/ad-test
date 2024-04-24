-- +goose Up
-- +goose StatementBegin
CREATE TABLE control_efficiency_statuses (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en VARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    color_code VARCHAR(25) NOT NULL,
    type VARCHAR(45) NOT NULL,
        CONSTRAINT CHK_control_implementation_statuses_type CHECK (type IN ('efficient', 'moderate', 'inefficient')),
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_control_implementation_statuses_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.control_efficiency_statuses;
-- +goose StatementEnd

