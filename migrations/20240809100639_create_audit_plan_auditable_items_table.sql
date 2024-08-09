-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_plan_auditable_items
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    audit_plan_id INT NOT NULL,
    auditable_item_id INT NOT NULL,
    user_id INT NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_audit_plan_auditable_items_audit_plans FOREIGN KEY (audit_plan_id) REFERENCES audit_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_auditable_items_auditable_items FOREIGN KEY (auditable_item_id) REFERENCES auditable_items(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_auditable_items_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,

);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE audit_plan_auditable_items;

-- +goose StatementEnd