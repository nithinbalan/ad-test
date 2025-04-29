-- +goose Up
-- +goose StatementBegin
CREATE TABLE grc_report_governance_summary (
    id INT PRIMARY KEY IDENTITY(1,1),
    grc_report_id INT NOT NULL,
    review_user_id INT NULL,
    approve_user_id INT NULL,
    policies_in_system NVARCHAR(MAX) NOT NULL,
    roles_defined INT NOT NULL DEFAULT 0,
    total_policies INT NOT NULL DEFAULT 0,
    policy_approved_rate DECIMAL(5,2) NULL,
    highlights NVARCHAR(MAX) NULL,
    challenges NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_grc_report_governance_summary_grc_reports FOREIGN KEY (grc_report_id) REFERENCES grc_reports(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_grc_report_governance_summary_user_1 FOREIGN KEY (review_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_governance_summary_user_2 FOREIGN KEY (approve_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_governance_summary_user_3 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_governance_summary_user_4 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE grc_report_governance_summary;
-- +goose StatementEnd
