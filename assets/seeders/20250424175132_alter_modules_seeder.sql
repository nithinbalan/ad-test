-- +goose Up
-- +goose StatementBegin
DECLARE @strategyModuleID INT;
DECLARE @swotModuleID INT;
DECLARE @pestleModuleID INT;
SELECT @strategyModuleID = id FROM modules WHERE identifier = 'strategy';
SELECT @swotModuleID = id FROM modules WHERE identifier = 'swot';
SELECT @pestleModuleID = id FROM modules WHERE identifier = 'pestle';

DELETE FROM role_permissions WHERE module_id = @swotModuleID;
DELETE FROM role_permissions WHERE module_id = @pestleModuleID;

DELETE FROM modules
WHERE identifier = 'swot';

DELETE FROM modules
WHERE identifier = 'pestle';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

INSERT INTO modules (identifier, title_en, title_ar, module_id, [order], icon, route)
VALUES 
    ('swot', 'SWOT', 'كدح', @strategyModuleID, 1, NULL, '/strategy/swot'),
    ('pestle', 'PESTLE', 'مدقة', @strategyModuleID, 2, NULL, '/strategy/pestle');
-- +goose StatementEnd
