-- +goose Up
-- +goose StatementBegin
-- Rename the table
EXEC sp_rename 'bcp_roles', 'bcp_teams';

-- Rename the constraints
EXEC sp_rename 'FK_bcp_roles_business_continuity_plans', 'FK_bcp_teams_business_continuity_plans', 'OBJECT';
EXEC sp_rename 'FK_bcp_roles_users_1', 'FK_bcp_teams_users_1', 'OBJECT';
EXEC sp_rename 'FK_bcp_roles_users_2', 'FK_bcp_teams_users_2', 'OBJECT';
EXEC sp_rename 'FK_bcp_roles_users_3', 'FK_bcp_teams_users_3', 'OBJECT';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Rename the table back to its original name
EXEC sp_rename 'bcp_teams', 'bcp_roles';

-- Rename the constraints back to their original names
EXEC sp_rename 'FK_bcp_teams_business_continuity_plans', 'FK_bcp_roles_business_continuity_plans', 'OBJECT';
EXEC sp_rename 'FK_bcp_teams_users_1', 'FK_bcp_roles_users_1', 'OBJECT';
EXEC sp_rename 'FK_bcp_teams_users_2', 'FK_bcp_roles_users_2', 'OBJECT';
EXEC sp_rename 'FK_bcp_teams_users_3', 'FK_bcp_roles_users_3', 'OBJECT';
-- +goose StatementEnd
