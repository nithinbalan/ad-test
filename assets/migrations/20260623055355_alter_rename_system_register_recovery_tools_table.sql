-- +goose Up
-- +goose StatementBegin
EXEC sp_rename 'system_register_recovery_tools', 'system_register_service_importances';
-- +goose StatementEnd

-- +goose StatementBegin
EXEC sp_rename 'CHK_system_register_recovery_tools_status', 'CHK_system_register_service_importances_status';
-- +goose StatementEnd

-- +goose StatementBegin
EXEC sp_rename 'FK_system_register_recovery_tools_users_1', 'FK_system_register_service_importances_users_1';
-- +goose StatementEnd

-- +goose StatementBegin
EXEC sp_rename 'FK_system_register_recovery_tools_users_2', 'FK_system_register_service_importances_users_2';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
EXEC sp_rename 'FK_system_register_service_importances_users_2', 'FK_system_register_recovery_tools_users_2';
-- +goose StatementEnd

-- +goose StatementBegin
EXEC sp_rename 'FK_system_register_service_importances_users_1', 'FK_system_register_recovery_tools_users_1';
-- +goose StatementEnd

-- +goose StatementBegin
EXEC sp_rename 'CHK_system_register_service_importances_status', 'CHK_system_register_recovery_tools_status';
-- +goose StatementEnd

-- +goose StatementBegin
EXEC sp_rename 'system_register_service_importances', 'system_register_recovery_tools';
-- +goose StatementEnd
