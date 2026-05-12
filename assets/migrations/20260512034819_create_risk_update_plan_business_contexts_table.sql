-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_update_plan_business_contexts (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_update_plan_timeline_id INT NOT NULL,
    internal_audit_note NVARCHAR(MAX) NULL,
    external_context NVARCHAR(MAX) NULL,
    annual_objective_description NVARCHAR(MAX) NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_risk_update_plan_business_contexts_risk_update_plan_timelines FOREIGN KEY (risk_update_plan_timeline_id) REFERENCES risk_update_plan_timelines(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_update_plan_business_contexts_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_update_plan_business_contexts;
-- +goose StatementEnd
