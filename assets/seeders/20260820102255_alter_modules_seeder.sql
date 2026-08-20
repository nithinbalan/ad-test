-- +goose Up
-- +goose StatementBegin
DECLARE @kpiModuleID INT;
DECLARE @kpiRegisterID INT;
DECLARE @kpiReviewID INT;
DECLARE @kpiDashboardID INT;
DECLARE @superAdminRoleID INT;

SELECT @kpiModuleID = id
FROM modules
WHERE identifier = 'kpi';

UPDATE modules
SET identifier = 'key_performance_indicator',
    title_en = 'KPI',
    title_ar = N'مؤشرات الأداء الرئيسية',
    icon = 'icon-kpi.svg',
    route = '/key-performance-indicator'
WHERE id = @kpiModuleID;

SELECT @kpiRegisterID = id
FROM modules
WHERE identifier = 'kpi_register';

UPDATE modules
SET title_en = 'KPI Register',
    title_ar = N'سجل مؤشرات الأداء الرئيسية',
    module_id = @kpiModuleID,
    [order] = 1,
    route = '/key-performance-indicator/kpi-register'
WHERE id = @kpiRegisterID;

INSERT INTO modules (
    identifier,
    title_en,
    title_ar,
    module_id,
    [order],
    icon,
    route
)
VALUES (
    'kpi_review',
    'KPI Reviews',
    N'مراجعات مؤشرات الأداء الرئيسية',
    @kpiModuleID,
    2,
    NULL,
    '/key-performance-indicator/kpi-review'
);

SELECT @kpiReviewID = SCOPE_IDENTITY();

UPDATE modules
SET [order] = [order] + 1
WHERE module_id = @kpiModuleID;

INSERT INTO modules (
    identifier,
    title_en,
    title_ar,
    module_id,
    [order],
    icon,
    route
)
VALUES (
    'kpi_dashboard',
    'Dashboard',
    N'لوحة القيادة',
    @kpiModuleID,
    1,
    NULL,
    '/key-performance-indicator/dashboard'
);

SELECT @kpiDashboardID = SCOPE_IDENTITY();

DECLARE @viewPermissionTypeID INT;
DECLARE @modifyPermissionTypeID INT;
DECLARE @managePermissionTypeID INT;
DECLARE @subModuleID INT;

SELECT @superAdminRoleID = id
FROM roles
WHERE is_super_admin = 1;

SELECT @viewPermissionTypeID = id
FROM role_permission_types
WHERE type = 'view';

SELECT @modifyPermissionTypeID = id
FROM role_permission_types
WHERE type = 'modify';

SELECT @managePermissionTypeID = id
FROM role_permission_types
WHERE type = 'manage';

DECLARE submodule_cursor CURSOR FOR
    SELECT id
    FROM modules
    WHERE module_id = @kpiModuleID;

OPEN submodule_cursor;
FETCH NEXT FROM submodule_cursor INTO @subModuleID;

WHILE @@FETCH_STATUS = 0
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM role_permissions
        WHERE role_id = @superAdminRoleID
          AND module_id = @subModuleID
          AND role_permission_type_id = @viewPermissionTypeID
    )
    BEGIN
        INSERT INTO role_permissions (
            role_id, module_id, role_permission_type_id,
            is_granted, created_by, created_at
        )
        VALUES (
            @superAdminRoleID, @subModuleID, @viewPermissionTypeID,
            1, @superAdminRoleID, GETDATE()
        );
    END;

    IF NOT EXISTS (
        SELECT 1
        FROM role_permissions
        WHERE role_id = @superAdminRoleID
          AND module_id = @subModuleID
          AND role_permission_type_id = @modifyPermissionTypeID
    )
    BEGIN
        INSERT INTO role_permissions (
            role_id, module_id, role_permission_type_id,
            is_granted, created_by, created_at
        )
        VALUES (
            @superAdminRoleID, @subModuleID, @modifyPermissionTypeID,
            1, @superAdminRoleID, GETDATE()
        );
    END;

    IF NOT EXISTS (
        SELECT 1
        FROM role_permissions
        WHERE role_id = @superAdminRoleID
          AND module_id = @subModuleID
          AND role_permission_type_id = @managePermissionTypeID
    )
    BEGIN
        INSERT INTO role_permissions (
            role_id, module_id, role_permission_type_id,
            is_granted, created_by, created_at
        )
        VALUES (
            @superAdminRoleID, @subModuleID, @managePermissionTypeID,
            1, @superAdminRoleID, GETDATE()
        );
    END;

    FETCH NEXT FROM submodule_cursor INTO @subModuleID;
END;

CLOSE submodule_cursor;
DEALLOCATE submodule_cursor;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @kpiReviewID INT;
DECLARE @kpiModuleID INT;
DECLARE @kpiDashboardID INT;
DECLARE @kpiDashboardOrder INT;

SELECT @kpiReviewID = id
FROM modules
WHERE identifier = 'kpi_review';

SELECT @kpiModuleID = id
FROM modules
WHERE identifier = 'key_performance_indicator';

SELECT
    @kpiDashboardID = id,
    @kpiDashboardOrder = [order]
FROM modules
WHERE identifier = 'kpi_dashboard';

IF @kpiDashboardID IS NOT NULL
BEGIN
    DELETE FROM role_permissions
    WHERE module_id = @kpiDashboardID;

    DELETE FROM activity_logs
    WHERE sub_module_id = @kpiDashboardID;

    DELETE FROM modules
    WHERE id = @kpiDashboardID;

    UPDATE modules
    SET [order] = [order] - 1
    WHERE module_id = @kpiModuleID
      AND [order] > @kpiDashboardOrder;
END;

-- Remove the KPI Reviews submodule
DELETE FROM role_permissions
WHERE module_id = @kpiReviewID;

DELETE FROM activity_logs
WHERE sub_module_id = @kpiReviewID;

DELETE FROM modules
WHERE id = @kpiReviewID;

-- Restore the original KPI routes and identifier
UPDATE modules
SET route = '/kpi/register'
WHERE identifier = 'kpi_register';

UPDATE modules
SET identifier = 'kpi',
    title_en = 'KPI',
    title_ar = N'مؤشرات الأداء الرئيسية',
    route = '/kpi'
WHERE id = @kpiModuleID;
-- +goose StatementEnd
