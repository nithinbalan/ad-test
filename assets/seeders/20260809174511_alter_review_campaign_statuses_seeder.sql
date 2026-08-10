-- +goose Up
-- +goose StatementBegin
INSERT INTO review_campaign_statuses (
    title_en,
    title_ar,
    color_code,
    type,
    status
)
VALUES (
    'Rejected',
    N'مرفوض',
    '#DC3545',
    'rejected',
    'active'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM review_campaign_statuses
WHERE type = 'rejected';
-- +goose StatementEnd
