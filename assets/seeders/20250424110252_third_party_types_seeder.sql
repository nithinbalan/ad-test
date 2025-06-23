-- +goose Up
-- +goose StatementBegin
INSERT INTO third_party_types (title_en, title_ar, status) VALUES
('Vendor', 'بائع', 'active'),
('Consultant', 'مستشار', 'active'),
('Service Provider', 'مزود خدمة', 'active'),
('Contractor', 'مقاول', 'active'),
('Supplier', 'مورد', 'active'),
('Partner', 'شريك', 'active'),
('Affiliate', 'شركة تابعة', 'active'),
('Distributor', 'موزع', 'active'),
('Reseller', 'بائع إعادة', 'active'),
('Outsourced Entity', 'جهة خارجية مستعان بها', 'active');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE third_party_types;
-- +goose StatementEnd
