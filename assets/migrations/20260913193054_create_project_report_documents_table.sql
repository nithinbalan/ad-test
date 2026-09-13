-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_report_documents (
    id INT IDENTITY(1, 1) PRIMARY KEY,
    project_report_project_id INT NOT NULL,
    [name] NVARCHAR(255) NULL,
    [size] INT NULL,
    [token] VARCHAR(MAX) NULL,
    original_name NVARCHAR(255) NULL,
    [type] VARCHAR(255) NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_project_report_documents_project_report_projects
        FOREIGN KEY (project_report_project_id)
        REFERENCES project_report_projects (id)
        ON DELETE CASCADE
        ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_report_documents;
-- +goose StatementEnd
