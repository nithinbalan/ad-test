-- +goose Up
-- +goose StatementBegin
UPDATE modules
SET route = '/management-review/meetings'
WHERE identifier = 'meeting';

UPDATE modules
SET route = '/management-review/action-plan-list'
WHERE identifier = 'meeting_action_plan';

UPDATE modules
SET route = '/management-review/report'
WHERE identifier = 'meeting_report';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE modules
SET route = '/review-meeting/meetings'
WHERE identifier = 'meeting';

UPDATE modules
SET route = '/review-meeting/action-plan-list'
WHERE identifier = 'meeting_action_plan';

UPDATE modules
SET route = '/review-meeting/report'
WHERE identifier = 'meeting_report';
-- +goose StatementEnd
