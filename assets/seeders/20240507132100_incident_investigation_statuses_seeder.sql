-- +goose Up
-- +goose StatementBegin
INSERT INTO incident_investigation_statuses
    (title_en, title_ar, color_code, type, status)
VALUES
    ('Open', 'مفتوح', '#FF0000', 'open', 'active'),
    ('Under Investigation', 'تحت التحقيق', '#FFFF00', 'under-investigation', 'active'),
    ('Closed', 'مغلق', '#808080', 'closed', 'active'),
    ('Resolved', 'محلول', '#00FF00', 'resolved', 'active'),
    ('On Hold', 'في الانتظار', '#0000FF', 'on-hold', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE incident_investigation_statuses;
-- +goose StatementEnd
