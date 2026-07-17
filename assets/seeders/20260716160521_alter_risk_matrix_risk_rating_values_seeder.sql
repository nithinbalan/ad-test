-- +goose Up
-- +goose StatementBegin

DECLARE
    @VeryLowID INT,
    @LowID INT,
    @MediumID INT,
    @HighID INT,
    @CatastrophicID INT;

SELECT @VeryLowID = id
FROM risk_matrix_risk_ratings
WHERE type = 'very_low'
AND risk_matrix_id = (SELECT id FROM risk_matrix WHERE type = 'isms');

SELECT @LowID = id
FROM risk_matrix_risk_ratings
WHERE type = 'low'
AND risk_matrix_id = (SELECT id FROM risk_matrix WHERE type = 'isms');

SELECT @MediumID = id
FROM risk_matrix_risk_ratings
WHERE type = 'medium'
AND risk_matrix_id = (SELECT id FROM risk_matrix WHERE type = 'isms');

SELECT @HighID = id
FROM risk_matrix_risk_ratings
WHERE type = 'high'
AND risk_matrix_id = (SELECT id FROM risk_matrix WHERE type = 'isms');

SELECT @CatastrophicID = id
FROM risk_matrix_risk_ratings
WHERE type = 'catastrophic'
AND risk_matrix_id = (SELECT id FROM risk_matrix WHERE type = 'isms');

INSERT INTO risk_matrix_risk_rating_values (risk_matrix_risk_rating_id, value)
VALUES
(@VeryLowID,1),
(@VeryLowID,2),

(@LowID,3),
(@LowID,4),
(@LowID,5),
(@LowID,6),

(@MediumID,8),
(@MediumID,9),
(@MediumID,10),
(@MediumID,12),

(@HighID,15),
(@HighID,16),

(@CatastrophicID,20),
(@CatastrophicID,25);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE v
FROM risk_matrix_risk_rating_values v
JOIN risk_matrix_risk_ratings r
ON v.risk_matrix_risk_rating_id = r.id
JOIN risk_matrix m
ON r.risk_matrix_id = m.id
WHERE m.type = 'isms';

-- +goose StatementEnd