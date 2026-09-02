-- +goose Up
-- +goose StatementBegin
CREATE TABLE penetration_test_retest_findings (
    id INT PRIMARY KEY IDENTITY(1, 1),
    penetration_test_retest_id INT NOT NULL,
    penetration_test_finding_id INT NOT NULL,
    is_verified BIT NOT NULL DEFAULT 0,
    verified_at DATETIME NULL,
    created_by_id INT NOT NULL,
    created_at DATETIME NOT NULL,
    CONSTRAINT UQ_penetration_test_retest_findings_retest_finding UNIQUE (penetration_test_retest_id, penetration_test_finding_id),
    CONSTRAINT FK_penetration_test_retest_findings_retests FOREIGN KEY (penetration_test_retest_id) REFERENCES penetration_test_retests (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_retest_findings_findings FOREIGN KEY (penetration_test_finding_id) REFERENCES penetration_test_findings (id) ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT FK_penetration_test_retest_findings_created_by FOREIGN KEY (created_by_id) REFERENCES users (id) ON DELETE NO ACTION ON UPDATE NO ACTION
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS penetration_test_retest_findings;
-- +goose StatementEnd
