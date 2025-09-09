-- +goose Up
-- +goose StatementBegin
ALTER TABLE third_party_assessments
ADD otp NVARCHAR(MAX) NULL;

ALTER TABLE third_party_assessments
ADD otp_created_at DATETIME NULL;

ALTER TABLE third_party_assessments
ADD otp_expiry DATETIME NULL;

ALTER TABLE third_party_assessments
ADD otp_attempts INT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE third_party_assessments
DROP COLUMN otp;

ALTER TABLE third_party_assessments
DROP COLUMN otp_created_at;

ALTER TABLE third_party_assessments
DROP COLUMN otp_expiry;

ALTER TABLE third_party_assessments
DROP COLUMN otp_attempts;
-- +goose StatementEnd
