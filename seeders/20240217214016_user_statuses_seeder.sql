-- +goose Up
-- +goose StatementBegin
INSERT INTO user_statuses
    (title_en, title_ar, color_code, type, status)
VALUES
    ('Active', 'نشيط', '#28A745', 'active', 'active'),
    ('On Leave', 'في إجازة', '#17A2B8', 'onleave', 'active'),
    ('Inactive', 'غير نشيط', '#6C757D', 'inactive', 'active'),
    ('Locked', 'مقفل', '#DC3545', 'locked', 'active'),
    ('Retired', 'متقاعد', '#FFA500', 'retired', 'active'),
    ('Resigned', 'مستقيل', '#D2B48C', 'resigned', 'active'),
    ('Deceased', 'متوفى', '#343A40', 'deceased', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE user_statuses;
-- +goose StatementEnd
