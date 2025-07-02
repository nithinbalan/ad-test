-- +goose Up
-- +goose StatementBegin
CREATE TABLE escalation_levels (
    id INT PRIMARY KEY IDENTITY(1, 1),
    escalation_id INT NOT NULL,
    day_offset INT NOT NULL,
    escalation_to VARCHAR(25) NOT NULL,
        CONSTRAINT CHK_escalation_levels_escalation_to CHECK (escalation_to IN ('reporting-user','designation')),
    notification_method VARCHAR(25) NOT NULL,
        CONSTRAINT CHK_escalation_notification_method CHECK (notification_method  IN ('app','email')),
    created_by_id INT,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_escalation_rule_users_1 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_escalation_rule_users_2 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_escalation_levels_escalations FOREIGN KEY (escalation_id) 
        REFERENCES escalations(id) ON DELETE CASCADE ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE dbo.escalation_levels;
-- +goose StatementEnd
