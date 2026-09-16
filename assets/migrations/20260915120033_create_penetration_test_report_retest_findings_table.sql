-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_report_retest_findings (
    id INT PRIMARY KEY IDENTITY(1,1),

    penetration_test_report_retest_id INT NOT NULL,
    penetration_test_report_finding_id INT NOT NULL,
    penetration_test_retest_finding_id INT NOT NULL,

    comment NVARCHAR(MAX) NULL,
    comment_at DATETIME NULL,
    result_title_en NVARCHAR(255) NULL,
    result_title_ar NVARCHAR(255) NULL,
    result_at DATETIME NULL,

    CONSTRAINT FK_penetration_test_report_retest_findings_retests FOREIGN KEY (penetration_test_report_retest_id) REFERENCES penetration_test_report_retests (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_report_retest_findings_findings FOREIGN KEY (penetration_test_report_finding_id) REFERENCES penetration_test_report_findings (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_report_retest_findings_source FOREIGN KEY (penetration_test_retest_finding_id) REFERENCES penetration_test_retest_findings (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE penetration_test_report_retest_findings;
-- +goose StatementEnd
