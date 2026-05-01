-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_response_workflows (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_response_id INT NOT NULL,
    level INT NOT NULL,
    workflow_level_type_id INT NOT NULL,
    approval_types VARCHAR(255) NOT NULL DEFAULT 'anyone',
    CONSTRAINT CHK_risk_response_wf_approval_types
        CHECK (approval_types IN ('anyone', 'all')),
    workflow_status_id INT NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_risk_response_wf_risk_response
        FOREIGN KEY (risk_response_id)
        REFERENCES risk_responses(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_risk_response_wf_workflow_level_types
        FOREIGN KEY (workflow_level_type_id)
        REFERENCES workflow_level_types(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_risk_response_wf_workflow_statuses
        FOREIGN KEY (workflow_status_id)
        REFERENCES workflow_statuses(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_risk_response_wf_users_1
        FOREIGN KEY (created_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION,

    CONSTRAINT FK_risk_response_wf_users_2
        FOREIGN KEY (updated_by_id)
        REFERENCES users(id)
        ON DELETE NO ACTION
        ON UPDATE NO ACTION
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.risk_response_workflows;
-- +goose StatementEnd
