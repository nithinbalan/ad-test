-- +goose Up
-- +goose StatementBegin
DECLARE @CriticalSeverityID INT;
DECLARE @HighSeverityID INT;
DECLARE @MediumSeverityID INT;
DECLARE @LowSeverityID INT;
DECLARE @SuperAdminRoleID INT;

SELECT @CriticalSeverityID = id FROM severity_levels WHERE type = 'critical';
SELECT @HighSeverityID = id FROM severity_levels WHERE type = 'high';
SELECT @MediumSeverityID = id FROM severity_levels WHERE type = 'medium';
SELECT @LowSeverityID = id FROM severity_levels WHERE type = 'low';
SELECT @SuperAdminRoleID = id FROM roles WHERE type = 'super-admin';

INSERT INTO sla_escalation_matrix (
    severity_id,
    remediation_days,
    is_email,
    is_alert,
    escalation_level_id,
    status
)
VALUES
    (@CriticalSeverityID, 3, 1, 1, @SuperAdminRoleID, 'active'),
    (@HighSeverityID, 7, 1, 0, @SuperAdminRoleID, 'active'),
    (@MediumSeverityID, 14, 0, 1, @SuperAdminRoleID, 'active'),
    (@LowSeverityID, 30, 0, 0, @SuperAdminRoleID, 'active');

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DELETE FROM sla_escalation_matrix;
-- +goose StatementEnd