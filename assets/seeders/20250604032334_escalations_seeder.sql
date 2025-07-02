-- +goose Up
-- +goose StatementBegin
DECLARE @complianceModuleID INT;

SET @complianceModuleID = (SELECT id FROM modules WHERE identifier = 'compliance');


INSERT INTO escalations (module_id,title_en, title_ar, status) VALUES
(@complianceModuleID,'Compliance Register Expiry', 'انتهاء صلاحية سجل الامتثال', 'active');

INSERT INTO escalations (module_id,title_en, title_ar, status) VALUES
(@complianceModuleID,'Compliance Requirement Deadline', 'الموعد النهائي لمتطلبات الامتثال', 'active');

INSERT INTO escalations (module_id,title_en, title_ar, status) VALUES
(@complianceModuleID,'Compliance Action Plan Overdue', 'خطة عمل الامتثال متأخرة', 'active');


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE escalations;
-- +goose StatementEnd