-- +goose Up
-- +goose StatementBegin
ALTER TABLE penetration_tests
DROP CONSTRAINT CHK_penetration_tests_status;

ALTER TABLE penetration_tests
DROP CONSTRAINT DF_penetration_tests_status;

ALTER TABLE penetration_tests
DROP COLUMN status;

ALTER TABLE penetration_tests
ADD penetration_test_status_id INT NULL,
    CONSTRAINT FK_penetration_tests_penetration_test_statuses
    FOREIGN KEY (penetration_test_status_id) REFERENCES penetration_test_statuses (id) ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE penetration_tests
DROP CONSTRAINT FK_penetration_tests_penetration_test_statuses;

ALTER TABLE penetration_tests
DROP COLUMN penetration_test_status_id;

ALTER TABLE penetration_tests
ADD status VARCHAR(20) NOT NULL
CONSTRAINT DF_penetration_tests_status DEFAULT 'draft'
CONSTRAINT CHK_penetration_tests_status CHECK (status IN ('draft', 'under retest', 'closed'));
-- +goose StatementEnd
