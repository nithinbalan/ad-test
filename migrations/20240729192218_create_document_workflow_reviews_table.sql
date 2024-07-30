-- +goose Up
-- +goose StatementBegin
CREATE TABLE document_workflow_reviews (
    id INT PRIMARY KEY IDENTITY(1,1),
    document_workflow_id INT NOT NULL,
    document_workflow_status_id INT NOT NULL,
    comment NVARCHAR(MAX),
    created_by INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by INT NOT NULL,
    updated_at DATETIME NOT NULL,
    CONSTRAINT FK_ocument_workflow_reviews_document_workflow FOREIGN KEY (document_workflow_id) REFERENCES document_workflows(id),
    CONSTRAINT FK_ocument_workflow_reviews_document_workflow_status FOREIGN KEY (document_workflow_status_id) REFERENCES document_workflow_statuses(id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.document_workflow_reviews;
-- +goose StatementEnd
