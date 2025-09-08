-- +goose Up
-- +goose StatementBegin

DECLARE @strategyID INT;
DECLARE @itemCount INT;

SELECT @strategyID = id FROM modules WHERE identifier = 'strategy';

SELECT @itemCount = COUNT(*) FROM kpi_library;

-- Update KPI Library module and item count
UPDATE knowledge_center_resources
SET module_id = @strategyID,
    item_count = @itemCount
WHERE title_en = 'KPI Library';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

-- Declare variables
DECLARE @riskID INT;

SELECT @riskID = id FROM modules WHERE identifier = 'risk';

UPDATE knowledge_center_resources
SET module_id = @riskID,
    item_count = 0
WHERE title_en = 'KPI Library';

-- +goose StatementEnd
