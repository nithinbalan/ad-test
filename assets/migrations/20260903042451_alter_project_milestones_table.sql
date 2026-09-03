-- +goose Up
-- +goose StatementBegin
ALTER TABLE project_milestones ADD project_phase_id INT NULL;

ALTER TABLE project_milestones ADD CONSTRAINT FK_project_milestones_project_phases
    FOREIGN KEY (project_phase_id) REFERENCES project_phases(id)
    ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE project_milestones DROP CONSTRAINT FK_project_milestones_project_phases;
ALTER TABLE project_milestones DROP COLUMN project_phase_id;
-- +goose StatementEnd
