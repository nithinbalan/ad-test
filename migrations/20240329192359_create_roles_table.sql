-- +goose Up
-- +goose StatementBegin
CREATE TABLE roles (
  id INT PRIMARY KEY IDENTITY(1, 1),
    title_en VARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
    CONSTRAINT CHK_roles_status CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.roles;
-- +goose StatementEnd
