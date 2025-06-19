-- +goose Up
-- +goose StatementBegin
UPDATE time_units SET 
    title_en = 'Minutes', 
    title_ar = 'دقائق', 
    type = 'minute' 
WHERE id = 1;

UPDATE time_units SET 
    title_en = 'Hours', 
    title_ar = 'ساعات', 
    type = 'hour' 
WHERE id = 2;

UPDATE time_units SET 
    title_en = 'Days', 
    title_ar = 'أيام', 
    type = 'day' 
WHERE id = 3;

UPDATE time_units SET 
    title_en = 'Weeks', 
    title_ar = 'أسابيع', 
    type = 'week' 
WHERE id = 4;
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
UPDATE time_units SET 
    title_en = 'Hours', 
    title_ar = '?????', 
    type = 'hour' 
WHERE id = 1;

UPDATE time_units SET 
    title_en = 'Minutes', 
    title_ar = '?????', 
    type = 'minute' 
WHERE id = 2;

UPDATE time_units SET 
    title_en = 'Days', 
    title_ar = '?????', 
    type = 'day' 
WHERE id = 3;

UPDATE time_units SET 
    title_en = 'Weeks', 
    title_ar = '??????', 
    type = 'week' 
WHERE id = 4;
-- +goose StatementEnd