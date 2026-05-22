-- +goose Up
-- +goose StatementBegin
BEGIN TRANSACTION;

DECLARE @RiskMatrixId INT;
SELECT @RiskMatrixId = id FROM risk_matrix WHERE type = 'erm';

-- Get rating IDs
DECLARE @LowId INT, @MediumId INT, @HighId INT, @CriticalId INT;

SELECT @LowId = id FROM risk_matrix_risk_ratings 
WHERE type = 'low' AND risk_matrix_id = @RiskMatrixId;

SELECT @MediumId = id FROM risk_matrix_risk_ratings 
WHERE type = 'medium' AND risk_matrix_id = @RiskMatrixId;

SELECT @HighId = id FROM risk_matrix_risk_ratings 
WHERE type = 'high' AND risk_matrix_id = @RiskMatrixId;

SELECT @CriticalId = id FROM risk_matrix_risk_ratings 
WHERE type = 'critical' AND risk_matrix_id = @RiskMatrixId;

INSERT INTO risk_matrix_risk_rating_values (risk_matrix_risk_rating_id, value)
VALUES 
(@LowId, 1),
(@LowId, 2),
(@LowId, 3);

INSERT INTO risk_matrix_risk_rating_values (risk_matrix_risk_rating_id, value)
VALUES 
(@MediumId, 4),
(@MediumId, 5),
(@MediumId, 6),
(@MediumId, 7),
(@MediumId, 8),
(@MediumId, 9);

INSERT INTO risk_matrix_risk_rating_values (risk_matrix_risk_rating_id, value)
VALUES 
(@HighId, 10),
(@HighId, 11),
(@HighId, 12),
(@HighId, 13),
(@HighId, 14),
(@HighId, 15);

INSERT INTO risk_matrix_risk_rating_values (risk_matrix_risk_rating_id, value)
VALUES 
(@CriticalId, 16),
(@CriticalId, 17),
(@CriticalId, 18),
(@CriticalId, 19),
(@CriticalId, 20),
(@CriticalId, 21),
(@CriticalId, 22),
(@CriticalId, 23),
(@CriticalId, 24),
(@CriticalId, 25);

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
BEGIN TRANSACTION;

DECLARE @RiskMatrixId INT;
SELECT @RiskMatrixId = id FROM risk_matrix WHERE type = 'erm';

DELETE FROM risk_matrix_risk_rating_values
WHERE risk_matrix_risk_rating_id IN (
    SELECT id FROM risk_matrix_risk_ratings 
    WHERE risk_matrix_id = @RiskMatrixId
);

COMMIT;
-- +goose StatementEnd