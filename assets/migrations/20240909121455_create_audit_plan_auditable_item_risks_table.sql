-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_plan_auditable_item_risks (
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_plan_auditable_item_id INT NOT NULL,
    risk_register_id INT NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_audit_plan_auditable_item_risks_audit_plan_auditable_items FOREIGN KEY (audit_plan_auditable_item_id) REFERENCES audit_plan_auditable_items(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_auditable_item_risks_risk_register FOREIGN KEY (risk_register_id) REFERENCES risk_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_auditable_item_risks_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_auditable_item_risks_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_plan_auditable_item_risks;
-- +goose StatementEnd
