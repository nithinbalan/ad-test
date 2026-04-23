-- +goose Up
-- +goose StatementBegin
CREATE TABLE kri_register_workflow_reviews (
    id INT PRIMARY KEY IDENTITY(1,1),
    kri_register_workflow_id INT NOT NULL,
    workflow_status_id INT NOT NULL,
    comment NVARCHAR(MAX),
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_krw_reviews_workflow
        FOREIGN KEY (kri_register_workflow_id)
        REFERENCES kri_register_workflows(id),

    CONSTRAINT FK_krw_reviews_workflow_status
        FOREIGN KEY (workflow_status_id)
        REFERENCES workflow_statuses(id)
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.kri_register_workflow_reviews;
-- +goose StatementEnd