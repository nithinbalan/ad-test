-- +goose Up
-- +goose StatementBegin
EXEC sp_rename 'awareness_campaigns.linked_sercurity_control', 'linked_security_control','COLUMN';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
EXEC sp_rename 'awareness_campaigns.linked_security_control', 'linked_sercurity_control','COLUMN';
-- +goose StatementEnd
