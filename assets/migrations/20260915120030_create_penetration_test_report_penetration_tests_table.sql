-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_report_penetration_tests (
    id INT PRIMARY KEY IDENTITY(1,1),

    penetration_test_report_id INT NOT NULL,
    penetration_test_id INT NOT NULL,
    reference_number NVARCHAR(50) NULL,

    system_application_name NVARCHAR(255) NOT NULL,
    scope NVARCHAR(MAX) NULL,
    year INT NOT NULL,
    period VARCHAR(2) NOT NULL,

    penetration_test_type_title_en NVARCHAR(255) NULL,
    penetration_test_type_title_ar NVARCHAR(255) NULL,
    penetration_test_plan_title_en NVARCHAR(255) NULL,
    penetration_test_plan_title_ar NVARCHAR(255) NULL,
    penetration_test_classification_title_en NVARCHAR(255) NULL,
    penetration_test_classification_title_ar NVARCHAR(255) NULL,
    penetration_test_executor_title_en NVARCHAR(255) NULL,
    penetration_test_executor_title_ar NVARCHAR(255) NULL,

    cs_responsible_employee_name_en NVARCHAR(255) NULL,
    cs_responsible_employee_name_ar NVARCHAR(255) NULL,

    penetration_test_status_title_en NVARCHAR(255) NULL,
    penetration_test_status_title_ar NVARCHAR(255) NULL,
    penetration_test_status_color_code NVARCHAR(255) NULL,

    environment_readiness_title_en NVARCHAR(255) NULL,
    environment_readiness_title_ar NVARCHAR(255) NULL,
    environment_notify_date DATETIME NULL,
    expected_ready_date DATETIME NULL,

    execution_start_date DATE NULL,
    execution_end_date DATE NULL,
    execution_status_title_en NVARCHAR(255) NULL,
    execution_status_title_ar NVARCHAR(255) NULL,
    security_result_title_en NVARCHAR(255) NULL,
    security_result_title_ar NVARCHAR(255) NULL,

    results_shared_with_management BIT NOT NULL DEFAULT 0,
    management_notify_date DATETIME NULL,
    management_response_date DATETIME NULL,
    remediation_notes NVARCHAR(MAX) NULL,

    closure_date DATE NULL,
    closure_remarks NVARCHAR(MAX) NULL,
    closure_owner_name_en NVARCHAR(255) NULL,
    closure_owner_name_ar NVARCHAR(255) NULL,

    created_at DATETIME NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_penetration_test_report_penetration_tests_penetration_test_reports FOREIGN KEY (penetration_test_report_id) REFERENCES penetration_test_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_report_penetration_tests_penetration_tests FOREIGN KEY (penetration_test_id) REFERENCES penetration_tests (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE penetration_test_report_penetration_tests;
-- +goose StatementEnd
