-- +goose Up
-- +goose StatementBegin

    DECLARE @organizationModuleID INT;

    SELECT @organizationModuleID = id FROM modules WHERE identifier = 'organization';

    UPDATE modules
    SET module_id = @organizationModuleID,
        title_en = N'Governance initiatives',
        title_ar = N'مبادرات الحوكمة'
    WHERE identifier = 'strategy_intiative';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DECLARE @strategyModuleID INT;

SELECT @strategyModuleID = id FROM modules WHERE identifier = 'strategy';

UPDATE modules
SET module_id = @strategyModuleID,
    title_en = N'Strategy Intiative',
    title_ar = N'المشاريع'
WHERE identifier = 'strategy_intiative';

-- +goose StatementEnd
