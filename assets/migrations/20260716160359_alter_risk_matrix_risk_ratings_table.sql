-- +goose Up
-- +goose StatementBegin

ALTER TABLE risk_matrix_risk_ratings
DROP CONSTRAINT CHK_risk_matrix_risk_ratings_type;

ALTER TABLE risk_matrix_risk_ratings
ALTER COLUMN type VARCHAR(20) NOT NULL;

ALTER TABLE risk_matrix_risk_ratings
ADD CONSTRAINT CHK_risk_matrix_risk_ratings_type
CHECK (
    type IN (
        'very_low',
        'low',
        'medium',
        'high',
        'very_high',
        'critical',
        'catastrophic'
    )
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE risk_matrix_risk_ratings
DROP CONSTRAINT CHK_risk_matrix_risk_ratings_type;

ALTER TABLE risk_matrix_risk_ratings
ALTER COLUMN type VARCHAR(10) NOT NULL;

ALTER TABLE risk_matrix_risk_ratings
ADD CONSTRAINT CHK_risk_matrix_risk_ratings_type
CHECK (
    type IN (
        'very_low',
        'low',
        'medium',
        'high',
        'very_high',
        'critical'
    )
);

-- +goose StatementEnd