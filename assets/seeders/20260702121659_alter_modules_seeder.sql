-- +goose Up
-- +goose StatementBegin
DECLARE @AwarenessReportingModuleID INT;
SELECT @AwarenessReportingModuleID = id FROM modules WHERE identifier = 'awareness_reporting';

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES ('awareness_post_campaign_survey', 'Post Campaign Survey', N'استبيان ما بعد الحملة', @AwarenessReportingModuleID, 3, NULL, '/awareness-reporting/post-campaign-survey');

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES ('awareness_report', 'Awareness Report', N'تقرير التوعية', @AwarenessReportingModuleID, 4, NULL, '/awareness-reporting/awareness-reports');

DECLARE @superAdminRoleID INT;
DECLARE @viewPermissionTypeID INT, @modifyPermissionTypeID INT, @managePermissionTypeID INT;
DECLARE @PostCampaignSurveyModuleID INT, @AwarenessReportModuleID INT;

SELECT @superAdminRoleID = id FROM roles WHERE is_super_admin = 1;
SELECT @viewPermissionTypeID = id FROM role_permission_types WHERE type = 'view';
SELECT @modifyPermissionTypeID = id FROM role_permission_types WHERE type = 'modify';
SELECT @managePermissionTypeID = id FROM role_permission_types WHERE type = 'manage';

SELECT @PostCampaignSurveyModuleID = id FROM modules WHERE identifier = 'awareness_post_campaign_survey';
SELECT @AwarenessReportModuleID = id FROM modules WHERE identifier = 'awareness_report';

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @PostCampaignSurveyModuleID 
      AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @PostCampaignSurveyModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @PostCampaignSurveyModuleID 
      AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @PostCampaignSurveyModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @PostCampaignSurveyModuleID 
      AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @PostCampaignSurveyModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @AwarenessReportModuleID 
      AND role_permission_type_id = @viewPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @AwarenessReportModuleID, @viewPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @AwarenessReportModuleID 
      AND role_permission_type_id = @modifyPermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @AwarenessReportModuleID, @modifyPermissionTypeID, 1, @superAdminRoleID, GETDATE());
END

IF NOT EXISTS (
    SELECT 1 FROM role_permissions 
    WHERE role_id = @superAdminRoleID 
      AND module_id = @AwarenessReportModuleID 
      AND role_permission_type_id = @managePermissionTypeID
)
BEGIN
    INSERT INTO role_permissions (role_id, module_id, role_permission_type_id, is_granted, created_by, created_at)
    VALUES (@superAdminRoleID, @AwarenessReportModuleID, @managePermissionTypeID, 1, @superAdminRoleID, GETDATE());
END
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM role_permissions
WHERE module_id IN (
    SELECT id FROM modules WHERE identifier IN ('awareness_post_campaign_survey', 'awareness_report')
);

DELETE FROM activity_logs
WHERE sub_module_id IN (
    SELECT id FROM modules WHERE identifier IN ('awareness_post_campaign_survey', 'awareness_report')
);

DELETE FROM modules
WHERE identifier IN ('awareness_post_campaign_survey', 'awareness_report');
-- +goose StatementEnd