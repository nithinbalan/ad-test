-- +goose Up
-- +goose StatementBegin

ALTER TABLE risk_evaluations
ADD 
    average_control_effectiveness_score DECIMAL(5,2) NULL;

-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

-- Then drop columns
ALTER TABLE risk_evaluations
DROP COLUMN 
    average_control_effectiveness_score;

-- +goose StatementEnd
