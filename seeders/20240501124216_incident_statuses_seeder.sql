-- +goose Up
-- +goose StatementBegin
INSERT INTO incident_statuses
    (title_en, title_ar, description_en, description_ar, color_code, type, status)
VALUES
    ('Open', 'مفتوح', 'The incident is open and under investigation.', 'الحادث مفتوح وقيد التحقيق.', '#FF0000', 'open', 'active'),
    ('Closed', 'مغلق', 'The incident has been closed after resolution and review.', 'تم إغلاق الحادث بعد حله ومراجعته.', '#808080', 'closed', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE incident_statuses;
-- +goose StatementEnd
