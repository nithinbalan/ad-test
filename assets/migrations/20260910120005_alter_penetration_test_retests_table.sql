-- +goose Up
-- +goose StatementBegin
ALTER TABLE penetration_test_retests
DROP CONSTRAINT CHK_penetration_test_retests_status;

ALTER TABLE penetration_test_retests
DROP COLUMN status;

ALTER TABLE penetration_test_retests
ADD penetration_test_retest_status_id INT NULL,
    CONSTRAINT FK_penetration_test_retests_penetration_test_retest_statuses
    FOREIGN KEY (penetration_test_retest_status_id) REFERENCES penetration_test_retest_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE penetration_test_retests
DROP CONSTRAINT FK_penetration_test_retests_penetration_test_retest_statuses;

ALTER TABLE penetration_test_retests
DROP COLUMN penetration_test_retest_status_id;

ALTER TABLE penetration_test_retests
ADD status VARCHAR(20) NOT NULL
CONSTRAINT DF_penetration_test_retests_status DEFAULT 'in-progress'
CONSTRAINT CHK_penetration_test_retests_status CHECK (status IN ('in-progress', 'completed'));
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE penetration_test_retests
DROP CONSTRAINT DF_penetration_test_retests_status;
-- +goose StatementEnd
