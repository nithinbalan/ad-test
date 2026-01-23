-- +goose Up
-- +goose StatementBegin


-- Create Consultation Request main module (NO permissions)
INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES
    ('consultation_request', 'Consultation Request', 'طلب استشارة', NULL, 15, 'icon-consultation.svg', '/consultations');


-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin


DELETE FROM modules
WHERE identifier = 'consultation_request';


-- +goose StatementEnd
