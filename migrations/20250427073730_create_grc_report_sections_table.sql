-- +goose Up
-- +goose StatementBegin
CREATE TABLE grc_report_sections (
    id INT PRIMARY KEY IDENTITY(1,1),
    grc_report_id INT NOT NULL,
    type VARCHAR(100) NOT NULL CHECK (type IN (
        'main-executive-summary',
        'governance-summary',
        'executive-summary-governance',
        'policy-approved-overview',
        'governance-tasks',
        'conclusion-governance',
        'executive-summary-compliance',
        'compliance-summary-self-assessment',
        'overall-compliance-status-by-department-self-assessment',
        'conclusion-recommendations-compliance',
        'integrated-grc-performance',
        'risk-management-summary',
        'executive-summary-risk',
        'risk-governance-context',
        'top-5-risks',
        'key-risk-indicators-kris-performance',
        'risk-control-matrix',
        'risk-monitoring-reporting-communication',
        'recommendations-action-plans',
        'conclusion-recommendations-risk',
        'main-conclusion'
    )),
    title_en NVARCHAR(255) NOT NULL,
    title_ar NVARCHAR(255) NOT NULL,
    description_en NVARCHAR(MAX) NULL,
    description_ar NVARCHAR(MAX) NULL,
    order_num INT NOT NULL,
    created_by_id INT NOT NULL,
    updated_by_id INT NULL,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NULL,
    CONSTRAINT FK_grc_report_sections_grc_reports FOREIGN KEY (grc_report_id) REFERENCES grc_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_sections_users_users_1 FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_grc_report_sections_users_users_2 FOREIGN KEY (updated_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE grc_report_sections;
-- +goose StatementEnd
