-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_plan_checklists(
    id INT PRIMARY KEY IDENTITY(1, 1),
    audit_plan_id INT NOT NULL,
    auditable_item_id INT NOT NULL,
    question NVARCHAR(255) NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_audit_plan_checklists_audit_plans FOREIGN KEY (audit_plan_id) REFERENCES audit_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_checklists_auditable_items FOREIGN KEY (auditable_item_id) REFERENCES auditable_items(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_checklists_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_checklists_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION

);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_plan_checklists;

-- +goose StatementEnd