-- +goose Up
-- +goose StatementBegin
DELETE FROM integrations WHERE title_en = 'Saudi National ID' AND integration_category_id = 1;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
INSERT INTO integrations (integration_category_id,title_en, title_ar, description_en, description_ar,logo_url,api_base_url,is_enable, created_at, created_by_id) VALUES
(1,'Saudi National ID', 'الهوية الوطنية السعودية', 'The Saudi National ID Card is proof of identity and citizenship for Saudi citizens.', 'بطاقة الهوية الوطنية السعودية هي إثبات الهوية والجنسية للمواطنين السعوديين.','Saudi National ID',NULL, 0,GETDATE(),1);
-- +goose StatementEnd
