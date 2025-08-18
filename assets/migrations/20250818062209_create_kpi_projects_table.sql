-- +goose Up
-- +goose StatementBegin
CREATE TABLE kpi_projects (
    id INT PRIMARY KEY IDENTITY(1, 1),
    kpi_register_id INT NOT NULL,
    project_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_kpi_projects_kpi_register FOREIGN KEY (kpi_register_id) REFERENCES kpi_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_projects_projects FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kpi_projects_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE kpi_projects;

-- +goose StatementEnd
