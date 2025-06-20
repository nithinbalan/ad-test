-- +goose Up
-- +goose StatementBegin
CREATE TABLE workflow_level_participants (
    id INT PRIMARY KEY IDENTITY(1,1),
    workflow_level_id INT NOT NULL,
    user_id INT NULL,
    role_id INT NULL,
    
    CONSTRAINT FK_workflow_level_participants_workflow_levels FOREIGN KEY (workflow_level_id) REFERENCES workflow_levels(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_workflow_level_participants_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_workflow_level_participants_roles FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE workflow_level_participants;
-- +goose StatementEnd
