-- +goose Up
-- +goose StatementBegin

ALTER TABLE risk_register
ADD
    is_cyber_risk BIT NOT NULL
        CONSTRAINT DF_risk_register_is_cyber_risk DEFAULT (0),
    justification NVARCHAR(MAX) NULL,
    confidentiality_score INT NULL,
    integrity_score INT NULL,
    availability_score INT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

ALTER TABLE risk_register
DROP CONSTRAINT DF_risk_register_is_cyber_risk;

ALTER TABLE risk_register
DROP COLUMN
    is_cyber_risk,
    justification,
    confidentiality_score,
    integrity_score,
    availability_score;

-- +goose StatementEnd