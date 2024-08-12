-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_plans
(
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_number VARCHAR(45) NOT NULL,
    audit_category_id INT NOT NULL,
    audit_team_id INT NULL,
    team_lead_id INT NULL,
    title NVARCHAR(255) NOT NULL,
    description NVARCHAR(MAX) NULL,
    start_date DATE NULL,
    target_date DATE NULL,
    purpose NVARCHAR(MAX) NULL,
    objective NVARCHAR(MAX) NULL,
    audit_status_id INT NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_audit_plans_audit_categories FOREIGN KEY (audit_category_id) REFERENCES audit_categories(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plans_audit_teams FOREIGN KEY (audit_team_id) REFERENCES audit_teams(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plans_users_1 FOREIGN KEY (team_lead_id) REFERENCES audit_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plans_audit_statuses FOREIGN KEY (audit_status_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plans_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plans_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.audit_plans;
-- +goose StatementEnd

