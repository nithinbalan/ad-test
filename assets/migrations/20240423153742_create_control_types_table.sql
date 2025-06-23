-- +goose Up
-- +goose StatementBegin
CREATE TABLE control_types (
    id INT PRIMARY KEY IDENTITY(1, 1),
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_control_type_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS control_types;
-- +goose StatementEnd
