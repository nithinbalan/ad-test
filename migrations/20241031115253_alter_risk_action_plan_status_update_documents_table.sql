-- +goose Up
-- +goose StatementBegin
ALTER TABLE risk_action_plan_status_update_documents
ADD original_name NVARCHAR(255) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_action_plan_status_update_documents
DROP COLUMN original_name;
-- +goose StatementEnd
