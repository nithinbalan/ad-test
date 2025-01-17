-- +goose Up
-- +goose StatementBegin
DECLARE @riskModuleID INT;

SELECT @riskModuleID = id FROM modules WHERE identifier = 'risk';

UPDATE modules
SET module_id = @riskModuleID
WHERE identifier = 'risk_report';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @assetRiskModuleID INT;

SELECT @assetRiskModuleID = id FROM modules WHERE identifier = 'asset_risk';

UPDATE modules
SET module_id = @assetRiskModuleID
WHERE identifier = 'risk_report';
-- +goose StatementEnd
