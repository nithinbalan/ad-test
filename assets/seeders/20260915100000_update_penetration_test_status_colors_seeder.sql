-- +goose Up
-- +goose StatementBegin
UPDATE penetration_test_statuses SET color_code = '#17A2B8' WHERE type = 'draft';
UPDATE penetration_test_statuses SET color_code = '#FF9F45' WHERE type = 'under retest';
UPDATE penetration_test_statuses SET color_code = '#DC3545' WHERE type = 'closed';
UPDATE penetration_test_retest_statuses SET color_code = '#17A2B8' WHERE type = 'draft';
UPDATE penetration_test_retest_statuses SET color_code = '#FF9F45' WHERE type = 'in-progress';
UPDATE penetration_test_retest_statuses SET color_code = '#8BC34A' WHERE type = 'completed';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE penetration_test_statuses SET color_code = CASE type WHEN 'draft' THEN '#9CA3AF' WHEN 'under retest' THEN '#F59E0B' WHEN 'closed' THEN '#10B981' ELSE color_code END;
UPDATE penetration_test_retest_statuses SET color_code = CASE type WHEN 'draft' THEN '#9CA3AF' WHEN 'in-progress' THEN '#F59E0B' WHEN 'completed' THEN '#10B981' ELSE color_code END;
-- +goose StatementEnd
