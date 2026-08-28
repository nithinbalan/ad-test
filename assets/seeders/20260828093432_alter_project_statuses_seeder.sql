-- +goose Up
-- +goose StatementBegin
INSERT INTO project_statuses (
    title_en, title_ar, description_en, description_ar, color_code, type, status
)
VALUES
    ('On Track', N'على المسار', 'Project progress is proceeding as planned.', N'تقدم المشروع يسير وفقاً للخطة.', '#28A745', 'on_track', 'active'),
    ('At Risk', N'معرض للخطر', 'Project progress is at risk and requires attention.', N'تقدم المشروع معرض للخطر ويتطلب اهتماماً.', '#FFC107', 'at_risk', 'active'),
    ('Delayed', N'متأخر', 'Project progress is behind the planned schedule.', N'تقدم المشروع متأخر عن الجدول الزمني المخطط.', '#DC3545', 'delayed', 'active');

DECLARE @onTrackStatusID INT;
DECLARE @atRiskStatusID INT;
DECLARE @delayedStatusID INT;

SELECT @onTrackStatusID = id
FROM project_statuses
WHERE type = 'on_track';

SELECT @atRiskStatusID = id
FROM project_statuses
WHERE type = 'at_risk';

SELECT @delayedStatusID = id
FROM project_statuses
WHERE type = 'delayed';

UPDATE p
SET project_status_id = CASE
    WHEN ps.type IN ('not_started', 'planning', 'in_progress', 'completed') THEN @onTrackStatusID
    WHEN ps.type = 'on_hold' THEN @atRiskStatusID
    WHEN ps.type = 'cancelled' THEN @delayedStatusID
END
FROM projects p
INNER JOIN project_statuses ps ON ps.id = p.project_status_id
WHERE ps.type IN ('not_started', 'planning', 'in_progress', 'on_hold', 'cancelled', 'completed');

DELETE FROM project_statuses
WHERE type IN ('not_started', 'planning', 'in_progress', 'on_hold', 'cancelled', 'completed');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
INSERT INTO project_statuses (
    title_en, title_ar, description_en, description_ar, color_code, type, status
)
VALUES
    ('Not Started', N'لم يبدأ', 'Project has been created but work has not yet begun.', N'تم إنشاء المشروع ولكن لم يبدأ العمل بعد.', '#6C757D', 'not_started', 'active'),
    ('Planning', N'تخطيط', 'Project is in the planning phase.', N'المشروع في مرحلة التخطيط.', '#17A2B8', 'planning', 'active'),
    ('In Progress', N'قيد التنفيذ', 'Project work is currently underway.', N'العمل على المشروع جارٍ حالياً.', '#007BFF', 'in_progress', 'active'),
    ('On Hold', N'معلق', 'Project is temporarily paused.', N'المشروع متوقف مؤقتاً.', '#FFC107', 'on_hold', 'active'),
    ('Cancelled', N'ملغى', 'Project has been cancelled.', N'تم إلغاء المشروع.', '#DC3545', 'cancelled', 'active'),
    ('Completed', N'مكتمل', 'Project has been successfully completed.', N'تم الانتهاء من المشروع بنجاح.', '#28A745', 'completed', 'active');

DECLARE @inProgressStatusID INT;
DECLARE @onHoldStatusID INT;
DECLARE @cancelledStatusID INT;

SELECT @inProgressStatusID = id
FROM project_statuses
WHERE type = 'in_progress';

SELECT @onHoldStatusID = id
FROM project_statuses
WHERE type = 'on_hold';

SELECT @cancelledStatusID = id
FROM project_statuses
WHERE type = 'cancelled';

UPDATE p
SET project_status_id = CASE
    WHEN ps.type = 'on_track' THEN @inProgressStatusID
    WHEN ps.type = 'at_risk' THEN @onHoldStatusID
    WHEN ps.type = 'delayed' THEN @cancelledStatusID
END
FROM projects p
INNER JOIN project_statuses ps ON ps.id = p.project_status_id
WHERE ps.type IN ('on_track', 'at_risk', 'delayed');

DELETE FROM project_statuses
WHERE type IN ('on_track', 'at_risk', 'delayed');
-- +goose StatementEnd
