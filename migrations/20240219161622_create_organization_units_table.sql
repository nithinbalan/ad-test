-- +goose Up
-- +goose StatementBegin
CREATE TABLE organization_units (
    id INT PRIMARY KEY IDENTITY(1,1),
    organization_level_id INT,
    organization_unit_id INT, -- Self-referencing foreign key
    title_en NVARCHAR(45) NOT NULL,
    title_ar NVARCHAR(45) NOT NULL,
    status VARCHAR(10) NOT NULL CHECK (status IN ('active', 'inactive')), 
    CONSTRAINT FK_organization_units_organization_levels FOREIGN KEY (organization_level_id) REFERENCES organization_levels(id),
    CONSTRAINT FK_organization_units_self FOREIGN KEY (organization_unit_id) REFERENCES organization_units(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.organization_units;
-- +goose StatementEnd
