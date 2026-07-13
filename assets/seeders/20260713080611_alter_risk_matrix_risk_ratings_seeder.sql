-- +goose Up
-- +goose StatementBegin

DECLARE @RiskMatrixID INT;

SELECT @RiskMatrixID = id
FROM risk_matrix
WHERE type = 'isms';

INSERT INTO risk_matrix_risk_ratings (
    risk_matrix_id,
    title_en,
    title_ar,
    type,
    color_code,
    status
)
VALUES
(@RiskMatrixID, 'Very Low', N'منخفض جدًا', 'very_low', '#28A745', 'active'),
(@RiskMatrixID, 'Low', N'منخفض', 'low', '#8BC34A', 'active'),
(@RiskMatrixID, 'Medium', N'متوسط', 'medium', '#FFC107', 'active'),
(@RiskMatrixID, 'High', N'عالي', 'high', '#FF9800', 'active'),
(@RiskMatrixID, 'Catastrophic', N'كارثي', 'catastrophic', '#FF0000', 'active');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM risk_matrix_risk_ratings
WHERE risk_matrix_id = (
    SELECT id FROM risk_matrix WHERE type = 'isms'
);

-- +goose StatementEnd