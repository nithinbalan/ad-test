-- +goose Up
-- +goose StatementBegin
DECLARE @bcmPolicyModuleID INT;

SELECT @bcmPolicyModuleID = id FROM modules WHERE identifier = 'bcm_policies';

DELETE FROM role_permissions WHERE module_id = @bcmPolicyModuleID;
DELETE FROM modules WHERE id = @bcmPolicyModuleID;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Rollback logic to re-insert bcm_policies
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES ('bcm_policies', 'BCM Policies', 'السياسات', (SELECT id FROM modules WHERE identifier = 'business_continuity'), 1, NULL, '/business-continuity/policies');
-- +goose StatementEnd
