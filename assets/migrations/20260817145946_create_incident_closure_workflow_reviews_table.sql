-- +goose Up
-- +goose StatementBegin
CREATE TABLE incident_closure_workflow_reviews (
    id INT PRIMARY KEY IDENTITY(1,1),
    incident_closure_workflow_id INT NOT NULL,
    workflow_status_id INT NOT NULL,
    comment NVARCHAR(MAX),
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_incident_closure_workflow_reviews_incident_closure_workflow FOREIGN KEY (incident_closure_workflow_id) REFERENCES incident_closure_workflows(id),
    CONSTRAINT FK_incident_closure_workflow_reviews_workflow_status FOREIGN KEY (workflow_status_id) REFERENCES workflow_statuses(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE incident_closure_workflow_reviews;
-- +goose StatementEnd
