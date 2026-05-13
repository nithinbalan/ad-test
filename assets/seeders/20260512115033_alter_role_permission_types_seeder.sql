-- +goose Up
-- +goose StatementBegin

INSERT INTO role_permission_types (title_en, title_ar, type, status)
VALUES (
    'Enable',
    'تمكين',
    'passive_enabled',
    'active'
);

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

DELETE FROM role_permission_types WHERE [type] = 'passive_enabled';
  
-- +goose StatementEnd