-- +goose Up
-- +goose StatementBegin
ALTER TABLE tprm_compliance_responses
ADD color_code NVARCHAR(25) NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tprm_compliance_responses
DROP COLUMN color_code;
-- +goose StatementEnd
