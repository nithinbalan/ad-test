-- +goose Up
-- +goose StatementBegin
CREATE TABLE organization_unit_attributes
(
    id INT PRIMARY KEY IDENTITY(1,1),
    department_id INT NOT NULL,
    title NVARCHAR(MAX) NULL,
    description NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_organization_unit_attributes_department FOREIGN KEY (department_id) REFERENCES organization_units(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_organization_unit_attributes_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_organization_unit_attributes_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE organization_unit_attributes;
-- +goose StatementEnd
