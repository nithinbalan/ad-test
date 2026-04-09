-- +goose Up
-- +goose StatementBegin

ALTER TABLE risk_impact_analyses
ADD is_selected BIT NOT NULL
    CONSTRAINT DF_risk_impact_analyses_is_selected DEFAULT 0;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin

ALTER TABLE risk_impact_analyses
DROP CONSTRAINT DF_risk_impact_analyses_is_selected;

ALTER TABLE risk_impact_analyses
DROP COLUMN is_selected;

-- +goose StatementEnd