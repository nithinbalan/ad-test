-- +goose Up
-- +goose StatementBegin
INSERT INTO user_statuses
    (title_en, title_ar, color_code, type, status, is_active)
VALUES
    ('Active', 'نشيط', '#28A745', 'active', 'active', 1),
    ('On Leave', 'في إجازة', '#17A2B8', 'onleave', 'active', 0),
    ('Locked', 'مقفل', '#DC3545', 'locked', 'active', 0),
    ('Retired', 'متقاعد', '#FFA500', 'retired', 'active', 0),
    ('Resigned', 'مستقيل', '#D2B48C', 'resigned', 'active', 0),
    ('Deceased', 'متوفى', '#343A40', 'deceased', 'active', 0); 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE user_statuses;
-- +goose StatementEnd
