-- +goose Up
-- +goose StatementBegin

INSERT INTO audit_finding_statuses (
    title_en,
    title_ar,
    color_code,
    type,
    is_closed
)
VALUES
(
    'Re-Assessment Pending',
    N'بانتظار إعادة التقييم',
    '#FFC107',
    're-assessment-pending',
    0
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM audit_finding_statuses
WHERE [type] IN (
    're-assessment-pending'
);

-- +goose StatementEnd