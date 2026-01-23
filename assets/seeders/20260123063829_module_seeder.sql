-- +goose Up
-- +goose StatementBegin


DECLARE @parentConsultationID INT;


-- Get MAIN consultation_request module ID (parent)
SELECT @parentConsultationID = id
FROM modules
WHERE identifier = 'consultation_request'
  AND module_id IS NULL;


-- Update CHILD module
UPDATE modules
SET
    module_id = @parentConsultationID,
    [order] = 1,
    identifier = 'consultation',
    route = '/consultation/requests'
WHERE identifier = 'consultation_request'
  AND module_id IS NOT NULL;


-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin


DECLARE @complianceModuleID INT;


-- Get Compliance module ID
SELECT @complianceModuleID = id
FROM modules
WHERE identifier = 'compliance';


-- Revert CHILD module
UPDATE modules
SET
    module_id = @complianceModuleID,
    identifier = 'consultation_request',
    [order] = 7,
    route = 'compliance/consultation'
WHERE identifier = 'consultation'
  AND module_id IS NOT NULL;


-- +goose StatementEnd
