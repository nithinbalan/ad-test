-- +goose Up
-- +goose StatementBegin
DECLARE @criticalSeverityID INT, @highSeverityID INT;

SELECT @criticalSeverityID = id FROM severity_levels WHERE type = 'critical';
SELECT @highSeverityID = id FROM severity_levels WHERE type = 'high';

UPDATE incident_sla_escalation_matrix
SET close_within_incident_owner = 0,
    close_within_response_team = 1
WHERE type = 'resolution_matrix'
  AND severity_id IN (@criticalSeverityID, @highSeverityID);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DECLARE @criticalSeverityID INT, @highSeverityID INT;

SELECT @criticalSeverityID = id FROM severity_levels WHERE type = 'critical';
SELECT @highSeverityID = id FROM severity_levels WHERE type = 'high';

UPDATE incident_sla_escalation_matrix
SET close_within_incident_owner = 1,
    close_within_response_team = 1
WHERE type = 'resolution_matrix'
  AND severity_id IN (@criticalSeverityID, @highSeverityID);
-- +goose StatementEnd
