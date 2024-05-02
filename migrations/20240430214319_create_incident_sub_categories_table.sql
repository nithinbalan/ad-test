-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_sub_categories (
    id INT PRIMARY KEY IDENTITY(1, 1),
    incident_category_id INT NOT NULL,
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    status VARCHAR(10) NOT NULL DEFAULT 'active',
        CONSTRAINT CHK_incident_sub_categories_status CHECK (status IN ('active', 'inactive')),
    CONSTRAINT FK_incident_sub_categories_categories FOREIGN KEY (incident_category_id) REFERENCES incident_categories (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.incident_sub_categories;
-- +goose StatementEnd
