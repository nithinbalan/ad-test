-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_closures (
    id INT PRIMARY KEY IDENTITY(1, 1),
    project_id INT NOT NULL,
    overall_project_outcome NVARCHAR(MAX) NULL,
    closure_reason NVARCHAR(MAX) NULL,
    closure_comments NVARCHAR(MAX) NULL,
    lesson_learned NVARCHAR(MAX) NULL,
    approver_id INT NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_project_closures_projects FOREIGN KEY (project_id)
        REFERENCES projects(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_closures_users_1 FOREIGN KEY (approver_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_closures_users_2 FOREIGN KEY (created_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_closures_users_3 FOREIGN KEY (updated_by_id)
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS project_closures;
-- +goose StatementEnd

