-- +goose Up
-- +goose StatementBegin
CREATE TABLE bcp_communication_plans (
    id INT PRIMARY KEY IDENTITY(1,1),
    business_continuity_plan_id INT NULL,
    bcp_phase_id INT NULL,
    title NVARCHAR(MAX) NULL,
    responsible_user_id INT NULL,
    responsible_unit_id INT NULL,
    recipient NVARCHAR(MAX) NULL,
    frequency NVARCHAR(MAX) NULL,
    channel NVARCHAR(MAX) NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
   CONSTRAINT FK_bcp_communication_plans_business_continuity_plans FOREIGN KEY (business_continuity_plan_id) REFERENCES business_continuity_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_bcp_communication_plans_bcp_phases FOREIGN KEY (bcp_phase_id) REFERENCES bcp_phases(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_bcp_communication_plans_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
   CONSTRAINT FK_bcp_communication_plans_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.bcp_communication_plans;
-- +goose StatementEnd
