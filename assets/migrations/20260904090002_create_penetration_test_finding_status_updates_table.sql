-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_finding_status_updates (
    id INT PRIMARY KEY IDENTITY(1, 1),
    penetration_test_finding_id INT NOT NULL,
    action_plan_status_id INT NOT NULL,
    percentage INT NOT NULL DEFAULT 0,
    comment NVARCHAR(MAX) NOT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    CONSTRAINT FK_penetration_test_finding_status_updates_penetration_test_finding
        FOREIGN KEY (penetration_test_finding_id) REFERENCES penetration_test_findings(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_finding_status_updates_action_plan_status
        FOREIGN KEY (action_plan_status_id) REFERENCES action_plan_statuses(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_finding_status_updates_created_by
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_finding_status_updates_updated_by
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE penetration_test_finding_status_updates;
-- +goose StatementEnd
