-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategic_objective_processes (
    id INT PRIMARY KEY IDENTITY(1, 1),
	strategic_objective_id INT NOT NULL,
    process_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
	CONSTRAINT FK_strategic_objective_processes_strategic_objectives FOREIGN KEY (strategic_objective_id) REFERENCES strategic_objectives(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategic_objective_processes_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	CONSTRAINT FK_strategic_objective_processes_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE strategic_objective_processes;

-- +goose StatementEnd
