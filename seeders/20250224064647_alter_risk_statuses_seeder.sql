-- +goose Up
-- +goose StatementBegin
INSERT INTO risk_statuses
    (title_en, title_ar, description_en, description_ar, color_code, type, status)
VALUES
    ('Submitted', 'تم التقديم', 'The risk has been submitted for review.', 'تم تقديم المخاطر للمراجعة.', '#1E90FF', 'submitted', 'active'),
    ('Reverted', 'تم الإرجاع', 'The risk has been reverted for further modifications.', 'تم إرجاع المخاطر لمزيد من التعديلات.', '#FF4500', 'reverted', 'active'),
    ('Rejected', 'مرفوض', 'The risk has been reviewed and rejected.', 'تم مراجعة المخاطر ورفضها.', '#DC143C', 'rejected', 'active'),
    ('Approved', 'مقبول', 'The risk has been reviewed and approved.', 'تم مراجعة المخاطر والموافقة عليها.', '#008000', 'approved', 'active'),
    ('Resubmitted', 'تم إعادة التقديم', 'The risk has been resubmitted after modifications.', 'تم إعادة تقديم المخاطر بعد التعديلات.', '#FFD700', 'resubmitted', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM risk_statuses WHERE type IN ('submitted', 'reverted', 'rejected', 'approved', 'resubmitted');
-- +goose StatementEnd
