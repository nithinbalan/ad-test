-- +goose Up
-- +goose StatementBegin

UPDATE modules
SET
    identifier = 'audit_insight_action_plans',
    is_passive = 1
WHERE identifier = 'action_plans'
  AND module_id = (
      SELECT id
      FROM modules
      WHERE identifier = 'audit_insights'
  );

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

UPDATE modules
SET
    identifier = 'action_plans',
    is_passive = 0
WHERE identifier = 'audit_insight_action_plans'
  AND module_id = (
      SELECT id
      FROM modules
      WHERE identifier = 'audit_insights'
  );

-- +goose StatementEnd