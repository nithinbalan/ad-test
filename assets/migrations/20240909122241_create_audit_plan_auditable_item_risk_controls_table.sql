-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_plan_auditable_item_risk_controls (
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_plan_auditable_item_risk_id INT NOT NULL,
    control_id INT NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,


    CONSTRAINT FK_audit_plan_auditable_item_risk_controls_audit_plan_auditable_item_risks FOREIGN KEY (audit_plan_auditable_item_risk_id) REFERENCES audit_plan_auditable_item_risks(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_auditable_item_risk_controls_controls FOREIGN KEY (control_id) REFERENCES controls(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_auditable_item_risk_controls_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_auditable_item_risk_controls_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_plan_auditable_item_risk_controls;
-- +goose StatementEnd
