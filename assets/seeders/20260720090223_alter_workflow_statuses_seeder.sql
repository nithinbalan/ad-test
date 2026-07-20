-- +goose Up
-- +goose StatementBegin

INSERT INTO workflow_statuses (
    title_en,
    title_ar,
    color_code,
    type,
    status
)
VALUES (
    'Resubmitted',
    N'تمت إعادة التقديم',
    '#007BFF',
    'resubmit',
    'active'
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM workflow_statuses
WHERE type = 'resubmit';

-- +goose StatementEnd