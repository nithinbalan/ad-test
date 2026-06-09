-- +goose Up
-- +goose StatementBegin

INSERT INTO internal_committee_statuses (
    title_en,
    title_ar,
    color_code,
    type,
    status
)
VALUES
(
    'Cancelled',
    'ملغى',
    '#6C757D',
    'cancelled',
    'active'
),
(
    'Completed',
    'مكتمل',
    '#198754',
    'completed',
    'active'
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DELETE FROM internal_committee_statuses
WHERE type IN ('cancelled', 'completed');

-- +goose StatementEnd