-- +goose Up
-- +goose StatementBegin
BEGIN TRANSACTION;

DECLARE @RiskMatrixId INT;
SELECT @RiskMatrixId = id FROM risk_matrix WHERE type = 'erm';

INSERT INTO risk_matrix_risk_ratings (risk_matrix_id, title_en, title_ar, type, color_code, status)
VALUES
    (@RiskMatrixId, 'Very Low', N'منخفض جداً', 'very_low', '#008000', 'active'), -- Green
    (@RiskMatrixId, 'Low', N'منخفض', 'low', '#FFFF00', 'active'), -- Yellow
    (@RiskMatrixId, 'Medium', N'متوسط', 'medium', '#FFA500', 'active'), -- Orange
    (@RiskMatrixId, 'High', N'عالي', 'high', '#FF0000', 'active'), -- Red
    (@RiskMatrixId, 'Very High', N'عالية جداً', 'very_high', '#B22222', 'active'); -- Brick Red

DECLARE @Counter INT = 1;

WHILE @Counter <= 25
BEGIN
    DECLARE @RatingType VARCHAR(20);
    SET @RatingType = CASE 
                         WHEN @Counter BETWEEN 1 AND 5 THEN 'very_low'
                         WHEN @Counter BETWEEN 6 AND 10 THEN 'low'
                         WHEN @Counter BETWEEN 11 AND 15 THEN 'medium'
                         WHEN @Counter BETWEEN 16 AND 20 THEN 'high'
                         WHEN @Counter BETWEEN 21 AND 25 THEN 'very_high'
                       END;

    INSERT INTO risk_matrix_risk_rating_values (risk_matrix_risk_rating_id, value)
    SELECT rr.id, @Counter AS value
    FROM risk_matrix_risk_ratings rr
    WHERE rr.risk_matrix_id = @RiskMatrixId
    AND rr.type = @RatingType;

    SET @Counter = @Counter + 1;
END;

COMMIT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE risk_matrix_risk_rating_values;
TRUNCATE TABLE risk_matrix_risk_ratings;
-- +goose StatementEnd
