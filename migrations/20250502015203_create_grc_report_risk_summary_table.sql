-- +goose Up
-- +goose StatementBegin
CREATE TABLE grc_report_risk_summary (
    id INT PRIMARY KEY IDENTITY(1,1),
    grc_report_id INT NOT NULL,
    review_user_id INT NULL,
    approve_user_id INT NULL,
    risk_identified_count INT NOT NULL DEFAULT 0,
    critical_risks NVARCHAR(MAX) NULL,
    treatment_strategies NVARCHAR(MAX) NULL,
    risk_score_summary NVARCHAR(MAX) NULL,
    kris NVARCHAR(MAX) NULL,
    risk_overview_chart NVARCHAR(MAX) NULL,
    risk_heatmap_chart NVARCHAR(MAX) NULL,
    top_risks NVARCHAR(MAX) NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    updated_by_id INT NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_grc_report_risk_summary_grc_reports FOREIGN KEY (grc_report_id) REFERENCES grc_reports(id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT FK_grc_report_risk_summary_users_1 FOREIGN KEY (review_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_risk_summary_users_2 FOREIGN KEY (approve_user_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_risk_summary_users_3 FOREIGN KEY (created_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_risk_summary_users_4 FOREIGN KEY (updated_by_id) REFERENCES users(id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE grc_report_risk_summary;
-- +goose StatementEnd
