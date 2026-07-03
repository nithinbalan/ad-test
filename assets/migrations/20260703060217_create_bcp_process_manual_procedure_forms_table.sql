-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_process_manual_procedure_forms (
    id INT IDENTITY(1,1) PRIMARY KEY,
    bcp_process_id INT NOT NULL,
    bcm_manual_procedure_form_id INT NOT NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_bcp_process_manual_procedure_forms_processes FOREIGN KEY (bcp_process_id) REFERENCES bcp_processes(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_manual_procedure_forms_bcm_manual_procedure_forms FOREIGN KEY (bcm_manual_procedure_form_id) REFERENCES bcm_manual_procedure_forms(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_process_manual_procedure_forms_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE bcp_process_manual_procedure_forms;
-- +goose StatementEnd
