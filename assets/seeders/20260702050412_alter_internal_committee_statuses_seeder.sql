-- +goose Up
-- +goose StatementBegin
UPDATE internal_committee_statuses
SET color_code = '#28A745'
WHERE type = 'completed';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE internal_committee_statuses
SET color_code = '#198754'
WHERE type = 'completed';
-- +goose StatementEnd
