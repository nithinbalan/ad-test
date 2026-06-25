-- +goose Up
-- +goose StatementBegin
CREATE TABLE process_execution_records (
    id INT IDENTITY(1,1) PRIMARY KEY,
    process_id INT NOT NULL,
    bcm_execution_records_id INT NOT NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_process_execution_records_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_process_execution_records_bcm_execution_records FOREIGN KEY (bcm_execution_records_id) REFERENCES bcm_execution_records(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_process_execution_records_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE process_execution_records;
-- +goose StatementEnd
