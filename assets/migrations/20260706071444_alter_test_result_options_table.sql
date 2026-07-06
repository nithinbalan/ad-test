-- +goose Up
-- +goose StatementBegin
ALTER TABLE test_result_options
DROP CONSTRAINT CHK_test_result_options_type;

ALTER TABLE test_result_options
ADD CONSTRAINT CHK_test_result_options_type
CHECK (
    type IN (
        'successful',
        'partially-successful',
        'failure',
        'not-tested'
    )
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE test_result_options
DROP CONSTRAINT CHK_test_result_options_type;

ALTER TABLE test_result_options
ADD CONSTRAINT CHK_test_result_options_type
CHECK (
    type IN (
        'successful',
        'partially-successful',
        'failure'
    )
);
-- +goose StatementEnd
