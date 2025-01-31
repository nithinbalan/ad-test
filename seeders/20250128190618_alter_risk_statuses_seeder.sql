-- +goose Up
-- +goose StatementBegin
INSERT INTO risk_statuses
    (title_en, title_ar, description_en, description_ar, color_code, type, status)
VALUES
    ('Identified', 'تم التعرف', 'The risk has been identified.', 'تم تحديد المخاطر', '#0000FF', 'identified', 'active')
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE risk_statuses;
-- +goose StatementEnd