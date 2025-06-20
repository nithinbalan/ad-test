-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_change_request_workflow_reviews (
    id INT PRIMARY KEY IDENTITY(1,1),
    document_change_request_workflow_id INT NOT NULL,
    document_change_request_workflow_status_id INT NOT NULL,
    comment NVARCHAR(MAX),
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_document_change_request_workflow_reviews_workflow FOREIGN KEY (document_change_request_workflow_id) REFERENCES document_change_request_workflows(id),
    CONSTRAINT FK_document_change_request_workflow_reviews_workflow_status FOREIGN KEY (document_change_request_workflow_status_id) REFERENCES document_change_request_workflow_statuses(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_change_request_workflow_reviews;
-- +goose StatementEnd
