-- +goose Up
-- +goose StatementBegin
 
INSERT INTO tprm_settings (is_export, is_external_vendor)
VALUES 
(1,0);  
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE tprm_settings;
-- +goose StatementEnd