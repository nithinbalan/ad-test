-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_findings (
    id INT PRIMARY KEY IDENTITY(1, 1),
    reference_number VARCHAR(45) NOT NULL UNIQUE,
    penetration_test_id INT NOT NULL,
    title NVARCHAR(255) NOT NULL,
    action_plan_status_id INT NOT NULL,
    severity_level_id INT NOT NULL,
    description NVARCHAR(MAX) NOT NULL,
    recommendation NVARCHAR(MAX) NOT NULL,
    owner_id INT NOT NULL,
    target_remediation_date DATE NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_penetration_test_findings_penetration_test
        FOREIGN KEY (penetration_test_id) REFERENCES penetration_tests(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_findings_severity_level
        FOREIGN KEY (severity_level_id) REFERENCES severity_levels(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_findings_action_plan_status
        FOREIGN KEY (action_plan_status_id) REFERENCES action_plan_statuses(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_findings_owner
        FOREIGN KEY (owner_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_findings_created_by
        FOREIGN KEY (created_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_findings_updated_by
        FOREIGN KEY (updated_by_id) REFERENCES users(id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE penetration_test_findings;
-- +goose StatementEnd
