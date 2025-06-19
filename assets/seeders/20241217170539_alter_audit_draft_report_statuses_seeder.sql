-- +goose Up
-- +goose StatementBegin
UPDATE audit_draft_report_statuses
SET color_code = '#28A745'
WHERE type = 'completed';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Revert the color_code back to its previous value (assuming it was '#FFC107')
UPDATE audit_draft_report_statuses
SET color_code = '#DC3545'
WHERE type = 'completed';
-- +goose StatementEnd