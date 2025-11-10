-- +goose Up
-- +goose StatementBegin
CREATE TABLE policy_report_key_metrics (
    id INT PRIMARY KEY IDENTITY(1,1),
    policy_report_id INT NOT NULL,
    clauses_defined_count INT NOT NULL DEFAULT 0,
    total_acknowledgement_count INT NOT NULL DEFAULT 0,
    acknowledgements_completed_count INT NOT NULL DEFAULT 0,
    acknowledgements_pending_count INT NOT NULL DEFAULT 0,
    acknowledgement_percentage NUMERIC(5,2) NOT NULL DEFAULT 0,
    risks_reported_count INT NOT NULL DEFAULT 0,
    associated_standards_count INT NOT NULL DEFAULT 0,
    associated_clauses_count INT NOT NULL DEFAULT 0,
    strategic_objectives_linked_count INT NOT NULL DEFAULT 0,
    approval_cycle_count INT NOT NULL DEFAULT 0,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,

    CONSTRAINT FK_policy_report_key_metrics_reports
        FOREIGN KEY (policy_report_id) REFERENCES policy_reports (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION,

    CONSTRAINT FK_policy_report_key_metrics_users_1
        FOREIGN KEY (created_by_id) REFERENCES users (id)
        ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE policy_report_key_metrics;
-- +goose StatementEnd
