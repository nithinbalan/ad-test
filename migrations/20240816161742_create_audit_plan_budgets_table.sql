-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_plan_budgets
(
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_plan_id INT NOT NULL,
    expense_category_id INT NOT NULL,
    title NVARCHAR(255) NOT NULL,
    estimated_cost NVARCHAR(45) NOT NULL,
    comments NVARCHAR(MAX) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_audit_plan_budgets_audit_plans FOREIGN KEY (audit_plan_id) REFERENCES audit_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_budgets_tool_types FOREIGN KEY (expense_category_id) REFERENCES expense_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_budgets_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_budgets_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.audit_plan_budgets;
-- +goose StatementEnd

