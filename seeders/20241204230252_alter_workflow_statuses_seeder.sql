-- +goose Up
-- +goose StatementBegin
UPDATE workflow_statuses
SET title_en = CASE 
                  WHEN type = 'submit' THEN 'Submitted'
                  WHEN type = 'revert' THEN 'Reverted'
                  WHEN type = 'reject' THEN 'Rejected'
                  WHEN type = 'approve' THEN 'Approved'
                  WHEN type = 'publish' THEN 'Published'
                  ELSE title_en -- Retain existing value if no match
               END,
    title_ar = CASE 
                  WHEN type = 'submit' THEN 'تم التقديم'
                  WHEN type = 'revert' THEN 'تمت الإعادة'
                  WHEN type = 'reject' THEN 'تم الرفض'
                  WHEN type = 'approve' THEN 'تمت الموافقة'
                  WHEN type = 'publish' THEN 'تم النشر'
                  ELSE title_ar -- Retain existing value if no match
               END;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE workflow_statuses
SET title_en = CASE 
                  WHEN type = 'submit' THEN 'Submitted'
                  WHEN type = 'revert' THEN 'Revert'
                  WHEN type = 'reject' THEN 'Reject'
                  WHEN type = 'approve' THEN 'Approve'
                  WHEN type = 'publish' THEN 'Published'
                  ELSE title_en -- Retain existing value if no match
               END,
    title_ar = CASE 
                  WHEN type = 'submit' THEN 'تم التقديم'
                  WHEN type = 'revert' THEN 'إعادة'
                  WHEN type = 'reject' THEN 'رفض'
                  WHEN type = 'approve' THEN 'موافقة'
                  WHEN type = 'publish' THEN 'تم النشر'
                  ELSE title_ar -- Retain existing value if no match
               END;
-- +goose StatementEnd
