-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_statuses
ADD score_from INT NULL;

ALTER TABLE compliance_statuses
ADD score_to INT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_statuses
DROP COLUMN score_from;

ALTER TABLE compliance_statuses
DROP COLUMN score_to;
-- +goose StatementEnd
