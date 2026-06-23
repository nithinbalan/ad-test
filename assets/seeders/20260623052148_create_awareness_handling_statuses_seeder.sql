-- +goose Up
-- +goose StatementBegin
INSERT INTO awareness_handling_statuses(title_en, title_ar,type)
VALUES
('New',N'جديد','new'),
('Received',N'تم الاستلام','received'),
('Under Investigation',N'قيد التحقيق','under-investigation'),
('Closed',N'مغلق','closed');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM awareness_handling_statuses
-- +goose StatementEnd
