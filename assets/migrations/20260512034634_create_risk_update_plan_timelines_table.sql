-- +goose Up
-- +goose StatementBegin

CREATE TABLE risk_update_plan_timelines
(
    id INT PRIMARY KEY IDENTITY(1, 1),
    risk_update_plan_id INT NOT NULL,
    risk_update_plan_timeline_type_id INT NOT NULL,
    is_business_context_updated BIT NOT NULL DEFAULT 0,
    is_not_applicable BIT NOT NULL DEFAULT 0,
    is_completed BIT NOT NULL DEFAULT 0,
    justification NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_risk_update_plan_timelines_risk_update_plans FOREIGN KEY (risk_update_plan_id) REFERENCES risk_register_update_plans (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_update_plan_timelines_risk_update_plan_timeline_types FOREIGN KEY (risk_update_plan_timeline_type_id) REFERENCES risk_update_plan_timeline_types (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_update_plan_timelines_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_update_plan_timelines_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_update_plan_timelines;
-- +goose StatementEnd
