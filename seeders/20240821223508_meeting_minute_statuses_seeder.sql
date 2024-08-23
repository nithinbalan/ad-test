-- +goose Up
INSERT INTO meeting_minute_statuses (title_en, title_ar, type, color_code)
VALUES 
    ('Informational', 'معلوماتي', 'informational', '#17A2B8'), -- Light Blue
    ('Open', 'مفتوح', 'open', '#DC3545'), -- Red
    ('Closed', 'مغلق', 'closed', '#28A745'); -- Green

-- +goose Down
TRUNCATE TABLE meeting_minute_statuses;
