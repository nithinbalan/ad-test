-- +goose Up
-- +goose StatementBegin

ALTER TABLE risk_control_analyses
ADD 
    control_effectiveness_score DECIMAL(5,2) NULL,
    control_effectiveness_rating_id INT NULL,
    risk_reduction DECIMAL(5,2) NULL;
-- Make control_efficiency_status_id nullable
ALTER TABLE risk_control_analyses
ALTER COLUMN control_efficiency_status_id INT NULL;

ALTER TABLE risk_control_analyses
ADD CONSTRAINT FK_risk_control_analyses_effectiveness_rating_id 
FOREIGN KEY (control_effectiveness_rating_id) REFERENCES control_effectiveness_ratings(id);
-- +goose StatementEnd



-- +goose Down
-- +goose StatementBegin

-- Drop FK first
ALTER TABLE risk_control_analyses
DROP CONSTRAINT FK_risk_control_analyses_effectiveness_rating_id;

-- Then drop columns
ALTER TABLE risk_control_analyses
DROP COLUMN 
    control_effectiveness_score,
    control_effectiveness_rating_id,
    risk_reduction;

-- Make control_efficiency_status_id non-nullable again
ALTER TABLE risk_control_analyses
ALTER COLUMN control_efficiency_status_id INT NULL;

-- +goose StatementEnd
