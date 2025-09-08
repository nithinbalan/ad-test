-- +goose Up
-- +goose StatementBegin
DECLARE @type_id INT;
SELECT @type_id = id FROM license_types WHERE type = 'active';

UPDATE roles
SET license_type_id = @type_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE roles
SET license_type_id = NULL;
-- +goose StatementEnd
