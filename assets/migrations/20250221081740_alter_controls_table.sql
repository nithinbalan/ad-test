-- +goose Up
-- +goose StatementBegin
ALTER TABLE controls 
ALTER COLUMN test_responsible_user_id INT NULL;

ALTER TABLE controls 
ALTER COLUMN test_frequency_id INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE controls 
ALTER COLUMN test_responsible_user_id INT NOT NULL;

ALTER TABLE controls 
ALTER COLUMN test_frequency_id INT NOT NULL;
-- +goose StatementEnd
