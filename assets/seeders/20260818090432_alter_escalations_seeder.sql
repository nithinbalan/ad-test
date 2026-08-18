-- +goose Up
-- +goose StatementBegin
DECLARE @identityAccessReviewModuleID INT;

SET @identityAccessReviewModuleID = (SELECT id FROM modules WHERE identifier = 'identity_access_review');

INSERT INTO escalations (module_id, title_en, title_ar, description_en, description_ar, status, [type]) VALUES
(
@identityAccessReviewModuleID, 
'Identity Review Campaign Approval Overdue',
N'تأخر اعتماد حملة مراجعة الهوية',
'Escalate when an identity review campaign is not approved by its due date.', 
N'تصعيد عند عدم اعتماد حملة مراجعة الهوية بحلول تاريخ الاستحقاق.', 
'active', 
'identity-review-campaign-approval-overdue'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM escalations WHERE [type] = 'identity-review-campaign-approval-overdue';
-- +goose StatementEnd
