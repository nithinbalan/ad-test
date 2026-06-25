-- +goose Up
-- +goose StatementBegin
CREATE TABLE process_manual_procedure_forms (
    id INT IDENTITY(1,1) PRIMARY KEY,
    process_id INT NOT NULL,
    bcm_manual_procedure_form_id INT NOT NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_process_manual_procedure_forms_processes FOREIGN KEY (process_id) REFERENCES processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_process_manual_procedure_forms_bcm_manual_procedure_forms FOREIGN KEY (bcm_manual_procedure_form_id) REFERENCES bcm_manual_procedure_forms(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_process_manual_procedure_forms_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE process_manual_procedure_forms;
-- +goose StatementEnd
