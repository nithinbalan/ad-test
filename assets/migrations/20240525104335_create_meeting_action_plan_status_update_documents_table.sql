-- +goose Up
-- +goose StatementBegin
CREATE TABLE meeting_action_plan_status_update_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    meeting_action_plan_status_update_id INT NOT NULL,
    url VARCHAR(1000) NOT NULL,
    name NVARCHAR(255) NOT NULL,
    size INT NOT NULL,
    token VARCHAR(255) NOT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) NOT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_meeting_action_plan_status_update_documents_meeting_action_plan_status_updates FOREIGN KEY (meeting_action_plan_status_update_id) REFERENCES meeting_action_plan_status_updates(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_action_plan_status_update_documents_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_meeting_action_plan_status_update_documents_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.meeting_action_plan_status_update_documents;
-- +goose StatementEnd
