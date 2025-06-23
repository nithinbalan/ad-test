-- +goose Up
-- +goose StatementBegin
CREATE TABLE reminder_rules (
    id INT PRIMARY KEY IDENTITY(1, 1),
    reminder_id INT NOT NULL,
    trigger_condition_id INT NOT NULL,
    title_en NVARCHAR(100) NOT NULL,
    title_ar NVARCHAR(100) NOT NULL,
    day_offset INT NOT NULL,
    type VARCHAR(25) NOT NULL,
        CONSTRAINT CHK_reminder_rules_type CHECK (type IN ('user-task')),
    notification_method VARCHAR(25) NOT NULL,
        CONSTRAINT CHK_reminder_notification_method CHECK (notification_method  IN ('app','email')),
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_reminder_rule_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_reminder_rule_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_reminder_rules_reminders FOREIGN KEY (reminder_id) 
        REFERENCES reminders(id) ON DELETE CASCADE ON UPDATE NO ACTION,
    CONSTRAINT FK_reminder_rules_trigger_conditions FOREIGN KEY (trigger_condition_id) 
        REFERENCES trigger_conditions(id) ON DELETE CASCADE ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.reminder_rules;
-- +goose StatementEnd
