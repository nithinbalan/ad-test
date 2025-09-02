-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_controls (
    id INT PRIMARY KEY IDENTITY(1, 1),
    control_id INT NOT NULL,
    project_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_project_control_controls FOREIGN KEY (control_id) REFERENCES controls(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_project_control_projects FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK__project_control_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE project_controls;
-- +goose StatementEnd