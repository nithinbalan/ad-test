-- +goose Up
-- +goose StatementBegin
ALTER TABLE penetration_test_retest_findings
DROP CONSTRAINT CHK_penetration_test_retest_findings_result;

ALTER TABLE penetration_test_retest_findings
DROP CONSTRAINT FK_penetration_test_retest_findings_created_by;

ALTER TABLE penetration_test_retest_findings
DROP COLUMN result, created_by_id, created_at;

ALTER TABLE penetration_test_retest_findings
ADD retest_result_id INT NULL,
    comment_by_id INT NULL,
    comment_at DATETIME NULL,
    result_by_id INT NULL,
    CONSTRAINT FK_penetration_test_retest_findings_penetration_test_retest_results
      FOREIGN KEY (retest_result_id) REFERENCES penetration_test_retest_results (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_retest_findings_comment_by
      FOREIGN KEY (comment_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_retest_findings_result_by
      FOREIGN KEY (result_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE penetration_test_retest_findings
DROP CONSTRAINT FK_penetration_test_retest_findings_penetration_test_retest_results;

ALTER TABLE penetration_test_retest_findings
DROP CONSTRAINT FK_penetration_test_retest_findings_comment_by;

ALTER TABLE penetration_test_retest_findings
DROP CONSTRAINT FK_penetration_test_retest_findings_result_by;

ALTER TABLE penetration_test_retest_findings
DROP COLUMN retest_result_id, comment_by_id, comment_at, result_by_id;

ALTER TABLE penetration_test_retest_findings
ADD result VARCHAR(20) NULL,
    created_by_id INT NOT NULL CONSTRAINT DF_tmp_ptrf_created_by DEFAULT 1,
    created_at DATETIME NOT NULL CONSTRAINT DF_tmp_ptrf_created_at DEFAULT SYSDATETIME(),
    CONSTRAINT CHK_penetration_test_retest_findings_result CHECK (result IN ('pass', 'fail')),
    CONSTRAINT FK_penetration_test_retest_findings_created_by
      FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE penetration_test_retest_findings
DROP CONSTRAINT DF_tmp_ptrf_created_by;

ALTER TABLE penetration_test_retest_findings
DROP CONSTRAINT DF_tmp_ptrf_created_at;
-- +goose StatementEnd
