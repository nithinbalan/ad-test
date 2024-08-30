-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_checklists
(
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_id INT NOT NULL,
    audit_plan_checklist_id INT NULL,
    question NVARCHAR(MAX) NOT NULL,
    conformance_id INT NOT NULL,
    comment NVARCHAR(MAX) NULL,
    answered_by_id INT NULL,
    answered_at DATETIME NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,
    status_id INT NOT NULL,
    CONSTRAINT FK_audit_checklists_audits FOREIGN KEY (audit_id) REFERENCES audits(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_checklists_audit_plan_checklists FOREIGN KEY (audit_plan_checklist_id) REFERENCES audit_plan_checklists(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_checklists_conformances FOREIGN KEY (conformance_id) REFERENCES conformances(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_checklists_users_1 FOREIGN KEY (answered_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_checklists_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_checklists_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_checklists_statuses FOREIGN KEY (status_id) REFERENCES statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_checklists;
-- +goose StatementEnd
