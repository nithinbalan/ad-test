-- +goose Up
-- +goose StatementBegin
CREATE TABLE job_team_members
(
    id INT PRIMARY KEY IDENTITY(1,1),
    job_team_id INT NOT NULL,
    designation_id INT NOT NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_job_team_members_job_teams FOREIGN KEY (job_team_id) REFERENCES job_teams(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_job_team_members_designation_1 FOREIGN KEY (designation_id) REFERENCES designations(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_job_team_members_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.job_team_members;
-- +goose StatementEnd