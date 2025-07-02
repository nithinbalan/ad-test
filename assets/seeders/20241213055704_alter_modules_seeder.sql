-- +goose Up
-- +goose StatementBegin
-- update business function route
UPDATE modules
SET route = '/business-continuity/business-functions'
WHERE identifier = 'business_functions';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE modules
SET route = '/business-continuity/functions'
WHERE identifier = 'business_functions';
-- +goose StatementEnd
