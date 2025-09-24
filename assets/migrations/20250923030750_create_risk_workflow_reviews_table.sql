-- +goose Up
-- +goose StatementBegin
CREATE TABLE risk_workflow_reviews (
    id INT PRIMARY KEY IDENTITY(1,1),
    risk_workflow_id INT NOT NULL,
    workflow_status_id INT NOT NULL,
    comment NVARCHAR(MAX),
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_risk_workflow_reviews_risk_workflow FOREIGN KEY (risk_workflow_id) REFERENCES risk_workflows(id),
    CONSTRAINT FK_risk_workflow_reviews_workflow_status FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.risk_workflow_reviews;
-- +goose StatementEnd
