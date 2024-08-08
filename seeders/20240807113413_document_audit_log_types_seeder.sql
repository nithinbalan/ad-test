-- +goose Up
-- +goose StatementBegin
INSERT INTO document_audit_log_types
    (title_en, title_ar, color_code, type)
VALUES
    ('Published', 'نُشر', '#28A745', 'published'),
    ('Viewed', 'تمت المشاهدة', '#007BFF', 'viewed'),
    ('Previewed', 'معاينة', '#6C757D', 'previewed'),
    ('Downloaded', 'تم التنزيل', '#FFC107', 'downloaded'),
    ('Acknowledged', 'تم التأكيد', '#DC3545', 'acknowledged');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE document_audit_log_types;
-- +goose StatementEnd
