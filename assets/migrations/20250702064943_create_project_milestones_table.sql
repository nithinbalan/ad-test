-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_milestones (
    id INT PRIMARY KEY IDENTITY(1,1),
    project_id INT NOT NULL,
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    planned_date DATE NULL,
    actual_date DATE NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    project_milestone_status_id INT NOT NULL,

    CONSTRAINT FK_project_milestones_projects FOREIGN KEY (project_id) 
        REFERENCES projects(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_milestones_users1 FOREIGN KEY (created_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_milestones_users2 FOREIGN KEY (updated_by_id) 
        REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_project_milestones_project_milestone_statuses FOREIGN KEY (project_milestone_status_id) 
        REFERENCES project_milestone_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.project_milestones;
-- +goose StatementEnd