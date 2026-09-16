-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_report_findings (
    id INT PRIMARY KEY IDENTITY(1,1),

    penetration_test_report_id INT NOT NULL,
    penetration_test_report_penetration_test_id INT NOT NULL,
    penetration_test_finding_id INT NOT NULL,
    reference_number NVARCHAR(50) NULL,

    title NVARCHAR(255) NOT NULL,
    description NVARCHAR(MAX) NULL,
    recommendation NVARCHAR(MAX) NULL,

    severity_level_title_en NVARCHAR(255) NULL,
    severity_level_title_ar NVARCHAR(255) NULL,
    severity_level_color_code NVARCHAR(255) NULL,

    owner_name_en NVARCHAR(255) NULL,
    owner_name_ar NVARCHAR(255) NULL,

    action_plan_status_title_en NVARCHAR(255) NULL,
    action_plan_status_title_ar NVARCHAR(255) NULL,
    action_plan_status_color_code NVARCHAR(255) NULL,

    percentage INT NULL,
    target_remediation_date DATE NULL,
    created_at DATETIME NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_penetration_test_report_findings_reports FOREIGN KEY (penetration_test_report_id) REFERENCES penetration_test_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_report_findings_penetration_tests FOREIGN KEY (penetration_test_report_penetration_test_id) REFERENCES penetration_test_report_penetration_tests (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_report_findings_source FOREIGN KEY (penetration_test_finding_id) REFERENCES penetration_test_findings (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE penetration_test_report_findings;
-- +goose StatementEnd
