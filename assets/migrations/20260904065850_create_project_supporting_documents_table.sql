-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_supporting_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    project_id INT NOT NULL,
    title_en NVARCHAR(255) NULL,
    title_ar NVARCHAR(255) NULL,
    created_by_id INT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    updated_at DATETIME NULL,

    CONSTRAINT FK_project_supporting_documents_projects FOREIGN KEY (project_id) 
        REFERENCES projects(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_supporting_documents_users_1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_supporting_documents_users_2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_supporting_documents;
-- +goose StatementEnd
