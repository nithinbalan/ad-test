-- +goose Up
-- +goose StatementBegin

ALTER TABLE third_party_register
ADD score FLOAT NULL;

ALTER TABLE third_party_register
ADD risk_level_id INT NULL;

ALTER TABLE third_party_register
ADD CONSTRAINT FK_third_party_register_risk_levels
FOREIGN KEY (risk_level_id) REFERENCES risk_levels(id)
ON DELETE NO ACTION ON UPDATE NO ACTION;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE third_party_register
DROP CONSTRAINT FK_third_party_register_risk_levels;


ALTER TABLE third_party_register
DROP COLUMN risk_level_id;

ALTER TABLE third_party_register
DROP COLUMN score;
-- +goose StatementEnd
