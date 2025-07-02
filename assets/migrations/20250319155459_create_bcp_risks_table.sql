-- +goose Up
-- +goose StatementBegin
CREATE TABLE dbo.bcp_risks (
    id INT IDENTITY(1,1) PRIMARY KEY,
    business_continuity_plan_id INT NOT NULL,
    risk_register_id INT NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_bcp_risks_business_continuity_plan FOREIGN KEY (business_continuity_plan_id) 
        REFERENCES dbo.business_continuity_plans(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_risks_risk_register FOREIGN KEY (risk_register_id) 
        REFERENCES dbo.risk_register(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_risks_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES dbo.users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_bcp_risks_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES dbo.users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS dbo.bcp_risks;
-- +goose StatementEnd
