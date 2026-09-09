-- +goose Up
-- +goose StatementBegin
ALTER TABLE penetration_tests
    ADD status VARCHAR(20) NOT NULL
    CONSTRAINT DF_penetration_tests_status DEFAULT 'draft'
    CONSTRAINT CHK_penetration_tests_status CHECK (status IN ('draft', 'under retest', 'closed'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE penetration_tests
    DROP CONSTRAINT CHK_penetration_tests_status;
ALTER TABLE penetration_tests
    DROP CONSTRAINT DF_penetration_tests_status;
ALTER TABLE penetration_tests
    DROP COLUMN status;
-- +goose StatementEnd
