-- +goose Up
-- +goose StatementBegin
CREATE TABLE organization_unit_business_contexts
(
    id INT PRIMARY KEY IDENTITY(1,1),
    department_id INT NOT NULL,
    internal_audit_note NVARCHAR(MAX) NULL,
    external_context NVARCHAR(MAX) NULL,
    other_field_header NVARCHAR(MAX) NULL,
    other_field_description NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_organization_unit_business_contexts_department FOREIGN KEY (department_id) REFERENCES organization_units(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_organization_unit_business_contexts_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_organization_unit_business_contexts_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE organization_unit_business_contexts;
-- +goose StatementEnd
