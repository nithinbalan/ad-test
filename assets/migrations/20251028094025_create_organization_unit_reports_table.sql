-- +goose Up
-- +goose StatementBegin
CREATE TABLE organization_unit_reports (
    id INT PRIMARY KEY IDENTITY(1,1),
    organization_unit_id INT NOT NULL,
    title NVARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_organization_unit_reports_organization_units FOREIGN KEY (organization_unit_id) REFERENCES organization_units (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_organization_unit_reports_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE organization_unit_reports;
-- +goose StatementEnd
