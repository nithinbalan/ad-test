-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_update_plan_risk_team_documents(
    id INT PRIMARY KEY IDENTITY(1,1),
     risk_update_plan_timeline_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token NVARCHAR(max) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME,
    CONSTRAINT FK_rup_risk_team_docs_timeline FOREIGN KEY (risk_update_plan_timeline_id) REFERENCES risk_update_plan_timelines(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_update_plan_risk_team_documents_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_risk_update_plan_risk_team_documents_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE risk_update_plan_risk_team_documents;
-- +goose StatementEnd
