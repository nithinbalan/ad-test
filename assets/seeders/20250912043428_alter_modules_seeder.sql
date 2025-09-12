-- +goose Up
-- +goose StatementBegin
UPDATE modules
SET route = '/settings/policy'
WHERE identifier = 'policy_management_setting';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE modules
SET route = '/setting/policy-setting'
WHERE identifier = 'policy_management_setting';
-- +goose StatementEnd
