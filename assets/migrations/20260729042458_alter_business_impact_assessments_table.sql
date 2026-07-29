-- +goose Up
-- +goose StatementBegin

ALTER TABLE business_impact_assessments
ADD
    is_selected BIT NOT NULL
        CONSTRAINT DF_business_impact_assessments_is_selected DEFAULT (0);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE business_impact_assessments
DROP CONSTRAINT DF_business_impact_assessments_is_selected;

ALTER TABLE business_impact_assessments
DROP COLUMN
    is_selected;
-- +goose StatementEnd