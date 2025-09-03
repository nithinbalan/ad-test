-- +goose Up
-- +goose StatementBegin
CREATE TABLE kri_processes (
    id INT PRIMARY KEY IDENTITY(1, 1),
    kri_register_id INT NOT NULL,
    process_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_kri_processes_kri_register FOREIGN KEY (kri_register_id) REFERENCES kri_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_kri_processes_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
	CONSTRAINT FK_kri_processes_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE kri_processes;

-- +goose StatementEnd
