-- +goose Up
-- +goose StatementBegin
INSERT INTO risk_statuses
    (title_en, title_ar, description_en, description_ar, color_code, type, status)
VALUES
    ('Open', 'مفتوح', 'The risk is open and has not been assessed yet.', 'المخاطر مفتوحة ولم يتم تقييمها بعد.', '#FF0000', 'open', 'active'),
    ('Assessed', 'تم التقييم', 'The risk has been assessed and its impact and likelihood have been determined.', 'تم تقييم المخاطر وتحديد تأثيرها واحتمالية حدوثها.', '#FFFF00', 'assessed', 'active'),
    ('Mitigated', 'تم التخفيف', 'Actions have been taken to mitigate the risk and reduce its impact.', 'تم اتخاذ إجراءات لتخفيف المخاطر وتقليل تأثيرها.', '#FFA500', 'mitigated', 'active'),
    ('Accepted', 'مقبول', 'The risk has been accepted and no further action is deemed necessary.', 'تم قبول المخاطر ولا يُعتبر اللازم اتخاذ إجراءات إضافية.', '#00FF00', 'accepted', 'active'),
    ('Closed', 'مغلق', 'The risk has been closed after assessment and mitigation.', 'تم إغلاق المخاطر بعد التقييم والتخفيف.', '#808080', 'closed', 'active'),
    ('Pending', 'معلق', 'The risk is pending further assessment or decision.', 'المخاطر معلقة بانتظار مزيد من التقييم أو القرار.', '#800080', 'pending', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE risk_statuses;
-- +goose StatementEnd
