-- +goose Up
-- +goose StatementBegin
CREATE TABLE project_action_plan_comments
(
    id INT PRIMARY KEY IDENTITY(1,1),
    project_action_plan_id INT NOT NULL,
    comment NVARCHAR(MAX) NOT NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL DEFAULT GETDATE(),
    CONSTRAINT FK_project_action_plan_comments_plan FOREIGN KEY (project_action_plan_id) REFERENCES project_action_plans(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_project_action_plan_comments_users FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE project_action_plan_comments;
-- +goose StatementEnd
