-- +goose Up
-- +goose StatementBegin
UPDATE grc_report_workflow_statuses
SET color_code = '#17A2B8'
WHERE type = 'draft';

UPDATE grc_report_workflow_statuses
SET color_code = '#FFA500'
WHERE type = 'in-review';

UPDATE grc_report_workflow_statuses
SET color_code = '#D6A2E8'
WHERE type = 'reviewed';

UPDATE grc_report_workflow_statuses
SET color_code = '#FD7E14'
WHERE type = 'pending-approval';

UPDATE grc_report_workflow_statuses
SET color_code = '#28A745'
WHERE type = 'approved';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE grc_report_workflow_statuses
SET color_code = '#0000FF'
WHERE type = 'draft';

UPDATE grc_report_workflow_statuses
SET color_code = '#FFD700'
WHERE type = 'in-review';

UPDATE grc_report_workflow_statuses
SET color_code = '#FF4500'
WHERE type = 'reviewed';

UPDATE grc_report_workflow_statuses
SET color_code = '#800080'
WHERE type = 'pending-approval';

UPDATE grc_report_workflow_statuses
SET color_code = '#008000'
WHERE type = 'approved';
-- +goose StatementEnd
