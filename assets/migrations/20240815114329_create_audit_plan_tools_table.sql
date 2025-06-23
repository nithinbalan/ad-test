-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_plan_tools
(
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_plan_id INT NOT NULL,
    tool_type_id INT NOT NULL,
    title NVARCHAR(255) NOT NULL,
    purpose NVARCHAR(MAX) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_audit_plan_tools_audit_plans FOREIGN KEY (audit_plan_id) REFERENCES audit_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_tools_tool_types FOREIGN KEY (tool_type_id) REFERENCES tool_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_tools_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_tools_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.audit_plan_tools;
-- +goose StatementEnd

