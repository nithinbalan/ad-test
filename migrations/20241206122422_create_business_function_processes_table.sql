-- +goose Up
-- +goose StatementBegin
CREATE TABLE business_function_processes (
    id INT PRIMARY KEY IDENTITY(1, 1),
    business_function_id INT NOT NULL,
    process_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_business_function_processes_business_functions FOREIGN KEY (business_function_id) REFERENCES business_functions(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_business_function_processes_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE business_function_processes;

-- +goose StatementEnd
