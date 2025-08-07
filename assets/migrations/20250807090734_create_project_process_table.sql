-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_processes (
    id INT PRIMARY KEY IDENTITY(1, 1),
    process_id INT NOT NULL,
    project_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_project_processes_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_project_processes_projects FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	CONSTRAINT FK_project_processes_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE project_processes;

-- +goose StatementEnd
