-- +goose Up
-- +goose StatementBegin
UPDATE lifecycle_statuses 
SET color_code = '#28A745' 
WHERE type = 'active';

UPDATE lifecycle_statuses 
SET color_code = '#DC3545' 
WHERE type = 'archived';

UPDATE lifecycle_statuses 
SET color_code = '#6C757D' 
WHERE type = 'expired';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE lifecycle_statuses 
SET color_code = '#00FF00' 
WHERE type = 'active';

UPDATE lifecycle_statuses 
SET color_code = '#808080' 
WHERE type = 'archived';

UPDATE lifecycle_statuses 
SET color_code = '#0000FF' 
WHERE type = 'expired';
-- +goose StatementEnd
