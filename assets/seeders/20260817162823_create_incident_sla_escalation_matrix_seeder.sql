-- +goose Up
-- +goose StatementBegin

DECLARE @criticalSeverityID INT, @highSeverityID INT, @mediumSeverityID INT, @lowSeverityID INT;
DECLARE @minutesUnitID INT, @hoursUnitID INT, @daysUnitID INT, @weeksUnitID INT;

SELECT @criticalSeverityID = id FROM severity_levels WHERE type = 'critical';
SELECT @highSeverityID = id FROM severity_levels WHERE type = 'high';
SELECT @mediumSeverityID = id FROM severity_levels WHERE type = 'medium';
SELECT @lowSeverityID = id FROM severity_levels WHERE type = 'low';

SELECT @minutesUnitID = id FROM business_continuity_duration_units WHERE type = 'minutes';
SELECT @hoursUnitID = id FROM business_continuity_duration_units WHERE type = 'hours';
SELECT @daysUnitID = id FROM business_continuity_duration_units WHERE type = 'days';
SELECT @weeksUnitID = id FROM business_continuity_duration_units WHERE type = 'weeks';

INSERT INTO incident_sla_escalation_matrix (
    severity_id, type, notify_within, notify_within_unit_id,
    notify_to_incident_owner, notify_to_response_team, notify_role_id,
    escalate_after_incident_owner, escalate_after_response_team, escalate_role_id,
    status
)
VALUES
(@criticalSeverityID, 'escalation_matrix', 15, @minutesUnitID, 0, 1, NULL, 0, 0, NULL, 'active'),
(@highSeverityID, 'escalation_matrix', 30, @minutesUnitID, 0, 1, NULL, 0, 0, NULL, 'active'),
(@mediumSeverityID, 'escalation_matrix', 2, @hoursUnitID, 0, 1, NULL, 0, 0, NULL, 'active'),
(@lowSeverityID, 'escalation_matrix', 1, @daysUnitID, 0, 1, NULL, 0, 0, NULL, 'active');

INSERT INTO incident_sla_escalation_matrix (
    severity_id, type, close_within, close_within_unit_id,
    close_within_incident_owner, close_within_response_team, close_within_role_id,
    status
)
VALUES
(@criticalSeverityID, 'resolution_matrix', 4, @hoursUnitID, 1, 1, NULL, 'active'),
(@highSeverityID, 'resolution_matrix', 1, @daysUnitID, 1, 1, NULL, 'active'),
(@mediumSeverityID, 'resolution_matrix', 3, @daysUnitID, 1, 0, NULL, 'active'),
(@lowSeverityID, 'resolution_matrix', 1, @weeksUnitID, 1, 0, NULL, 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM incident_sla_escalation_matrix;

-- +goose StatementEnd
