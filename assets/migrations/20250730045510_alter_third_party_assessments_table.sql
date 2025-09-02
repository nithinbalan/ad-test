-- +goose Up
-- +goose StatementBegin
ALTER TABLE third_party_assessments
ADD link NVARCHAR(MAX) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE third_party_assessments
DROP COLUMN link;
-- +goose StatementEnd
