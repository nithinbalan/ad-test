-- +goose Up
-- +goose StatementBegin

UPDATE control_types SET type = 'preventive' WHERE title_en = 'Preventive';
UPDATE control_types SET type = 'corrective' WHERE title_en = 'Corrective';
UPDATE control_types SET type = 'detective'  WHERE title_en = 'Detective';
UPDATE control_types SET type = 'deterrent'  WHERE title_en = 'Deterrent';

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

UPDATE control_types SET type = NULL
WHERE title_en IN ('Preventive', 'Corrective', 'Detective', 'Deterrent');

-- +goose StatementEnd