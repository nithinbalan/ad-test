-- +goose Up
-- +goose StatementBegin
ALTER TABLE business_continuity_plans
ADD is_test_and_exercise_completed BIT NOT NULL CONSTRAINT DF_business_continuity_plans_is_test_and_exercise_completed DEFAULT 0;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE business_continuity_plans
DROP CONSTRAINT DF_business_continuity_plans_is_test_and_exercise_completed;
ALTER TABLE business_continuity_plans
DROP COLUMN is_test_and_exercise_completed;
-- +goose StatementEnd
