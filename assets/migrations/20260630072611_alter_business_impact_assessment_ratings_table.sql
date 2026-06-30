-- +goose Up
-- +goose StatementBegin
ALTER TABLE business_impact_assessment_ratings
ADD override_reason NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE business_impact_assessment_ratings
DROP COLUMN override_reason;
-- +goose StatementEnd
