-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_report_retests (
    id INT PRIMARY KEY IDENTITY(1,1),

    penetration_test_report_id INT NOT NULL,
    penetration_test_report_penetration_test_id INT NOT NULL,
    penetration_test_retest_id INT NOT NULL,
    reference_number NVARCHAR(50) NULL,

    start_date DATE NOT NULL,
    end_date DATE NULL,
    executor_title_en NVARCHAR(255) NULL,
    executor_title_ar NVARCHAR(255) NULL,
    status_title_en NVARCHAR(255) NULL,
    status_title_ar NVARCHAR(255) NULL,
    status_color_code NVARCHAR(255) NULL,
    notes NVARCHAR(MAX) NULL,
    created_at DATETIME NULL,
    updated_at DATETIME NULL,

    CONSTRAINT FK_penetration_test_report_retests_reports FOREIGN KEY (penetration_test_report_id) REFERENCES penetration_test_reports (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_report_retests_penetration_tests FOREIGN KEY (penetration_test_report_penetration_test_id) REFERENCES penetration_test_report_penetration_tests (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_report_retests_source FOREIGN KEY (penetration_test_retest_id) REFERENCES penetration_test_retests (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE penetration_test_report_retests;
-- +goose StatementEnd
