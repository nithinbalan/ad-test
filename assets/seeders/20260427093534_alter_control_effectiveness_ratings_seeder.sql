-- +goose Up
-- +goose StatementBegin
UPDATE control_effectiveness_ratings SET color_code = '#FF0000' WHERE type = 'control_missing';
UPDATE control_effectiveness_ratings SET color_code = '#DC3545' WHERE type = 'ineffective';
UPDATE control_effectiveness_ratings SET color_code = '#FFC107' WHERE type = 'improvement_needed';
UPDATE control_effectiveness_ratings SET color_code = '#90DD75' WHERE type = 'partially_effective';
UPDATE control_effectiveness_ratings SET color_code = '#28A745' WHERE type = 'very_effective';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
UPDATE control_effectiveness_ratings SET color_code = NULL WHERE type = 'control_missing';
UPDATE control_effectiveness_ratings SET color_code = NULL WHERE type = 'ineffective';
UPDATE control_effectiveness_ratings SET color_code = NULL WHERE type = 'improvement_needed';
UPDATE control_effectiveness_ratings SET color_code = NULL WHERE type = 'partially_effective';
UPDATE control_effectiveness_ratings SET color_code = NULL WHERE type = 'very_effective';
-- +goose StatementEnd
