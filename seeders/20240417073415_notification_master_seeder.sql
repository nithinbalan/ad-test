-- +goose Up
-- +goose StatementBegin
INSERT INTO notification_masters (module_id, identifier, title_en, title_ar)
VALUES
    (3, 'new_user', 'New User Registration', 'تسجيل مستخدم جديد'),  
    (3, 'user_update', 'User Profile Update', 'تحديث ملف المستخدم'),
    (3, 'user_update', 'Reset Password', 'تحديث ملف المستخدم'),
    (2, 'new_post', 'Strategic Profile Update', 'إنشاء منشور جديد'),
    (4, 'post_update', 'New Process', 'تحديث المنشور');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE notification_masters;
-- +goose StatementEnd
