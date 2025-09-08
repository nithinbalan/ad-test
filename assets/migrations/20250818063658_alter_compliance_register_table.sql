-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_register
ADD score DECIMAL(10,2) NULL;
-- +goose StatementEnd
 
-- +goose Down
-- +goose StatementBegin
 
ALTER TABLE compliance_register
DROP COLUMN score;
-- +goose StatementEnd