-- +goose Up
-- +goose StatementBegin
ALTER TABLE job_teams
ALTER COLUMN team_lead_designation_id INT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE job_teams
ADD team_lead_user_id INT NULL
CONSTRAINT FK_job_teams_users_3 FOREIGN KEY REFERENCES users(id) 
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE job_teams
DROP CONSTRAINT FK_job_teams_users_3;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE job_teams
DROP COLUMN team_lead_user_id;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE job_teams
ALTER COLUMN team_lead_designation_id INT NOT NULL;
-- +goose StatementEnd
