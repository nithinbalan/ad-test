-- +goose Up
-- +goose StatementBegin
ALTER TABLE compliance_importance_classifications
ADD CONSTRAINT CHK_compliance_importance_classifications_type_values
CHECK (type IN (
    'normal',
    'urgent',
    'very_urgent',
    'important_and_very_urgent'
));
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE compliance_importance_classifications
DROP CONSTRAINT CHK_compliance_importance_classifications_type_values;
-- +goose StatementEnd
