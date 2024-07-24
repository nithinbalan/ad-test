-- +goose Up
-- +goose StatementBegin
CREATE TABLE audit_team_members
(
    id INT PRIMARY KEY IDENTITY(1,1),
    audit_team_id INT NOT NULL,
    user_id INT NOT NULL,
    created_by_id INT,
    created_at DATETIME NOT NULL,
    CONSTRAINT FK_audit_team_members_audit_teams FOREIGN KEY (audit_team_id) REFERENCES audit_teams(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_team_members_users_1 FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_audit_team_members_users_2 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.audit_team_members;
-- +goose StatementEnd
