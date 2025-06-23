-- +goose Up
-- +goose StatementBegin
UPDATE grc_report_comment_statuses
SET color_code = '#17A2B8'
WHERE type = 'new';

UPDATE grc_report_comment_statuses
SET color_code = '#FFA500'
WHERE type = 'resolved';

UPDATE grc_report_comment_statuses
SET color_code = '#28A745'
WHERE type = 'closed';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE grc_report_comment_statuses
SET color_code = '#0000FF'
WHERE type = 'new';

UPDATE grc_report_comment_statuses
SET color_code = '#FFD700'
WHERE type = 'resolved';

UPDATE grc_report_comment_statuses
SET color_code = '#008000'
WHERE type = 'closed';
-- +goose StatementEnd
