-- +goose Up
-- +goose StatementBegin
ALTER TABLE projects ADD budget DECIMAL(18, 2) NULL;
ALTER TABLE projects ADD scope NVARCHAR(MAX) NULL;
ALTER TABLE projects ADD project_phase_id INT NULL;
ALTER TABLE projects ADD project_classification_id INT NULL;

ALTER TABLE projects ADD CONSTRAINT FK_projects_project_phases
    FOREIGN KEY (project_phase_id) REFERENCES project_phases(id)
    ON DELETE SET NULL ON UPDATE NO ACTION;

ALTER TABLE projects ADD CONSTRAINT FK_projects_project_classifications
    FOREIGN KEY (project_classification_id) REFERENCES project_classifications(id)
    ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE projects DROP CONSTRAINT FK_projects_project_classifications;
ALTER TABLE projects DROP CONSTRAINT FK_projects_project_phases;
ALTER TABLE projects DROP COLUMN project_classification_id;
ALTER TABLE projects DROP COLUMN project_phase_id;
ALTER TABLE projects DROP COLUMN scope;
ALTER TABLE projects DROP COLUMN budget;
-- +goose StatementEnd