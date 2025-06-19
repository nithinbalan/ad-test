-- +goose Up
-- +goose StatementBegin
ALTER TABLE risk_statuses DROP CONSTRAINT CHK_risk_statuses_type;

ALTER TABLE risk_statuses ADD CONSTRAINT CHK_risk_statuses_type 
CHECK (type IN ('open', 'assessed', 'mitigated', 'accepted', 'closed', 'pending', 'identified', 'submitted', 'reverted', 'rejected', 'approved', 'resubmitted'));
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE risk_statuses DROP CONSTRAINT CHK_risk_statuses_type;

ALTER TABLE risk_statuses ADD CONSTRAINT CHK_risk_statuses_type 
CHECK (type IN ('identified','open', 'assessed', 'mitigated', 'accepted', 'closed', 'pending'));
-- +goose StatementEnd
