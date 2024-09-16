-- +goose Up
-- +goose StatementBegin
CREATE TABLE role_permission_types (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en VARCHAR(100) NOT NULL,
    title_ar VARCHAR(100) NOT NULL,
    type VARCHAR(20) CHECK (type IN ('manage', 'modify', 'view')) NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'inactive'))
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE role_permission_types;
-- +goose StatementEnd
