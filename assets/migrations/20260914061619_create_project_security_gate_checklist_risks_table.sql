-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_security_gate_checklist_risks (
    id INT PRIMARY KEY IDENTITY(1,1),
    checklist_id INT NOT NULL,
    risk_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_psg_checklist_risks_checklist
        FOREIGN KEY (checklist_id) REFERENCES project_security_gate_checklists(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_psg_checklist_risks_risk_register
        FOREIGN KEY (risk_register_id) REFERENCES risk_register(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_psg_checklist_risks_created_by
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_security_gate_checklist_risks;
-- +goose StatementEnd
