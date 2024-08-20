-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_plan_stakeholder_engagements
(
    id INT PRIMARY KEY IDENTITY(1,1),
    reference_number VARCHAR(45) NOT NULL UNIQUE,
    audit_plan_id INT NOT NULL,
    audit_engagement_type_id INT NOT NULL,
    title NVARCHAR(255) NOT NULL,
    objectives NVARCHAR(MAX) NULL,
    responsible_user_id INT NOT NULL,
    audit_schedule_status_id INT,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_audit_plan_stakeholder_engagements_engagement_types FOREIGN KEY (audit_engagement_type_id) REFERENCES audit_engagement_types(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_stakeholder_engagements_audit_plans FOREIGN KEY (audit_plan_id) REFERENCES audit_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_stakeholder_engagements_users_1 FOREIGN KEY (responsible_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_stakeholder_engagements_audit_schedule_statuses FOREIGN KEY (audit_schedule_status_id) REFERENCES audit_schedule_statuses(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_stakeholder_engagements_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_plan_stakeholder_engagements_users_3 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.audit_plan_stakeholder_engagements;
-- +goose StatementEnd
