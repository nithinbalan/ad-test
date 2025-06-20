-- +goose Up
-- +goose StatementBegin
CREATE TABLE strategy_profile_action_plan_status_update_documents (
    id INT PRIMARY KEY IDENTITY(1,1),
    strategy_profile_action_plan_status_update_id INT NOT NULL,
    url VARCHAR(1000) DEFAULT NULL,
    name NVARCHAR(255) NULL,
    size INT DEFAULT NULL,
    token VARCHAR(255) DEFAULT NULL,
    original_name NVARCHAR(255) NULL,
    type VARCHAR(255) DEFAULT NULL,
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_strategy_profile_action_plan_status_update_documents FOREIGN KEY (strategy_profile_action_plan_status_update_id) REFERENCES strategy_profile_action_plan_status_updates(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_action_plan_status_update_documents_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_strategy_profile_action_plan_status_update_documents_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE strategy_profile_action_plan_status_update_documents;
-- +goose StatementEnd
