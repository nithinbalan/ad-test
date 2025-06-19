-- +goose Up
-- +goose StatementBegin
INSERT INTO bcp_phases (title_en, title_ar, time_frame, status) VALUES
('Immediate', 'مباشر', '0-24 hrs', 'active'),
('Short-Term', 'على المدى القصير', '24-72 hrs', 'active'),
('Mid-Term', 'منتصف المدة', '72 hrs - 2 weeks', 'active'),
('Long-Term', 'طويل الأمد', '2+ weeks', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE bcp_phases;
-- +goose StatementEnd
