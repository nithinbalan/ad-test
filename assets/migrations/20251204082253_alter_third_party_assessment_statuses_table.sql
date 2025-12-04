-- +goose Up
-- +goose StatementBegin
ALTER TABLE third_party_assessment_statuses
DROP CONSTRAINT CHK_third_party_assessment_statuses_type;

ALTER TABLE third_party_assessment_statuses
ADD CONSTRAINT CHK_third_party_assessment_statuses_type 
CHECK (type IN ('draft','new','in-progress', 'submitted', 'approved', 'rejected'))
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE third_party_assessment_statuses
DROP CONSTRAINT CHK_third_party_assessment_statuses_type;

ALTER TABLE third_party_assessment_statuses
ADD CONSTRAINT CHK_third_party_assessment_statuses_type 
CHECK (type IN ('draft','new','in-progress', 'submitted'))
-- +goose StatementEnd
