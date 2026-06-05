-- +goose Up
-- +goose StatementBegin
CREATE TABLE workflow_custom_roles (
    id INT PRIMARY KEY IDENTITY(1,1),
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    module_identifier NVARCHAR(100) NOT NULL,
    CONSTRAINT CHK_workflow_custom_roles_type
        CHECK (type IN ('department_head'))
);
-- +goose StatementEnd




-- +goose Down
-- +goose StatementBegin
DROP TABLE workflow_custom_roles;
-- +goose StatementEnd

